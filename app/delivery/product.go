package delivery

import (
	"net/http"
	"strconv"

	"Tokobelanja/app/delivery/middleware"
	"Tokobelanja/app/helper"
	"Tokobelanja/domain"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type ProductHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(r *gin.RouterGroup, productUsecase domain.ProductUsecase) {
	handler := &ProductHandler{productUsecase}
	productRoute := r.Group("/products")
	productRoute.Use(middleware.Authentication())
	productRoute.GET("/", handler.GetProducts)
	productRoute.Use(middleware.Authorization([]string{"admin"}))
	productRoute.POST("/", handler.StoreProduct)
	productRoute.PUT(":productId", handler.UpdateProduct)
	productRoute.DELETE(":productId", handler.DeleteProduct)
}

func (c *ProductHandler) GetProducts(ctx *gin.Context) {
	products, err := c.productUsecase.GetProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": products})
}

func (c *ProductHandler) StoreProduct(ctx *gin.Context) {
	type StoreProduct struct {
		Title      string `json:"title" validate:"required"`
		Price      int64  `json:"price" validate:"required,gte=0,lte=50000000"`
		Stock      int64  `json:"stock" validate:"required,gte=5"`
		CategoryID int64  `json:"category_id" validate:"required"`
	}
	var storeProduct StoreProduct
	err := ctx.ShouldBindJSON(&storeProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = helper.ValidateStruct(storeProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var product domain.Product
	copier.Copy(&product, &storeProduct)
	productData, err := c.productUsecase.StoreProduct(ctx.Request.Context(), &product)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"data": gin.H{
			"id":          productData.ID,
			"title":       productData.Title,
			"price":       productData.Price,
			"stock":       productData.Stock,
			"category_Id": productData.CategoryID,
			"created_at":  productData.CreatedAt,
		},
	})
}

func (c *ProductHandler) UpdateProduct(ctx *gin.Context) {
	type UpdateProduct struct {
		Title      string `json:"title" validate:"required"`
		Price      int64  `json:"price" validate:"required,gte=0,lte=50000000"`
		Stock      int64  `json:"stock" validate:"required,gte=5"`
		CategoryID int64  `json:"category_id" validate:"required"`
	}
	var updateProduct UpdateProduct
	err := ctx.ShouldBindJSON(&updateProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = helper.ValidateStruct(updateProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var product domain.Product
	copier.Copy(&product, &updateProduct)
	productId, _ := strconv.ParseInt(ctx.Param("productId"), 10, 64)
	product.ID = productId
	productData, err := c.productUsecase.UpdateProduct(ctx, &product)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"id":          productData.ID,
			"title":       productData.Title,
			"price":       productData.Price,
			"stock":       productData.Stock,
			"category_Id": productData.CategoryID,
			"created_at":  productData.CreatedAt,
			"updated_at":  productData.UpdatedAt,
		},
	})

}

func (c *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productId, _ := strconv.ParseInt(ctx.Param("productId"), 10, 64)
	err := c.productUsecase.DeleteProduct(ctx, productId)
	if err != nil {
		ctx.JSON(getStatusCode(err), gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Product has been successfully deleted"})
}
