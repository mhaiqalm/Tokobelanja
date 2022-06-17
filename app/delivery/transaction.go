package delivery

import (
	"net/http"

	"Tokobelanja/app/delivery/middleware"
	"Tokobelanja/app/helper"
	"Tokobelanja/domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type TransactionHandler struct {
	transactionUsecase domain.TransactionUsecase
}

func NewTransactionHandler(r *gin.RouterGroup, transactionUsecase domain.TransactionUsecase) {
	handler := &TransactionHandler{transactionUsecase}
	transactionRoute := r.Group("/transactions")
	transactionRoute.Use(middleware.Authentication())
	transactionRoute.POST("/", handler.StoreTransaction)
	transactionRoute.GET("/my-transactions", handler.GetMyTransactions)

	transactionRoute.Use(middleware.Authorization([]string{"admin"}))
	transactionRoute.GET("/user-transactions", handler.GetTransactions)
}

func (t *TransactionHandler) GetTransactions(ctx *gin.Context) {
	transactions, err := t.transactionUsecase.GetTransactions(ctx)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": transactions})
}

func (t *TransactionHandler) GetMyTransactions(ctx *gin.Context) {
	userAuth := ctx.MustGet("user").(jwt.MapClaims)
	userID := int64(userAuth["id"].(float64))
	transactions, err := t.transactionUsecase.GetMyTransactions(ctx, userID)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": transactions})
}

func (t *TransactionHandler) StoreTransaction(ctx *gin.Context) {
	type StoreTransactionVal struct {
		ProductID int64 `json:"product_id" validate:"required"`
		Quantity  int64 `json:"quantity" validate:"required"`
	}
	var storeProductVal StoreTransactionVal
	err := ctx.ShouldBindJSON(&storeProductVal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = helper.ValidateStruct(storeProductVal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var transaction domain.Transaction
	copier.Copy(&transaction, &storeProductVal)

	userAuth := ctx.MustGet("user").(jwt.MapClaims)
	userID := int64(userAuth["id"].(float64))

	transaction.UserID = userID
	dataTransaction, err := t.transactionUsecase.StoreTransaction(ctx, &transaction)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	newDataTransaction := dataTransaction.(map[string]interface{})
	ctx.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "You have successfully purchased the product",
		"data": gin.H{
			"total_price":   newDataTransaction["total_price"],
			"quantity":      newDataTransaction["quantity"],
			"product_title": newDataTransaction["product_title"],
		},
	})
}
