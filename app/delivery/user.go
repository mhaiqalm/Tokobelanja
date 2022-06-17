package delivery

import (
	"net/http"
	"strconv"

	"Tokobelanja/app/delivery/middleware"
	"Tokobelanja/app/helper"
	"Tokobelanja/domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.RouterGroup, userUsecase domain.UserUsecase) {
	handler := &UserHandler{userUsecase}
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})
	userRoute := r.Group("/users")
	userRoute.POST("/register", handler.Register)
	userRoute.POST("/login", handler.Login)
	userRoute.Use(middleware.Authentication())
	userRoute.PATCH("/topup", handler.TopUp)
}

func (u *UserHandler) Register(ctx *gin.Context) {
	type UserRegister struct {
		FullName string `json:"full_name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=16"`
	}
	var userRegister UserRegister
	err := ctx.ShouldBindJSON(&userRegister)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = helper.ValidateStruct(userRegister)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var user domain.User
	copier.Copy(&user, &userRegister)
	userData, err := u.userUsecase.Register(ctx.Request.Context(), &user)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"data": gin.H{
			"id":         userData.ID,
			"full_name":  userData.FullName,
			"email":      userData.Email,
			"created_at": userData.CreatedAt,
		},
	})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type UserLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	var userLogin UserLogin
	err := ctx.ShouldBindJSON(&userLogin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = helper.ValidateStruct(userLogin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var user domain.User
	copier.Copy(&user, &userLogin)
	token, err := u.userUsecase.Login(ctx.Request.Context(), &user)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"token": token,
		},
	})

}

func (u *UserHandler) TopUp(ctx *gin.Context) {
	type Request struct {
		Balance int64 `json:"balance" validate:"required,gte=0,lte=100000000"`
	}
	var request Request

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = helper.ValidateStruct(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user domain.User
	copier.Copy(&user, &request)
	userAuth := ctx.MustGet("user").(jwt.MapClaims)
	userID := int64(userAuth["id"].(float64))
	user.ID = userID

	userData, err := u.userUsecase.TopUp(ctx.Request.Context(), &user)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"data": gin.H{
			"message": "Your balance has been successfully update to Rp" + strconv.Itoa(int(userData.Balance)),
		},
	})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrBalanceNotEnough:
		return http.StatusPaymentRequired
	case domain.ErrStockNotEnough:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
