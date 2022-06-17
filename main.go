package main

import (
	"os"

	_handler "Tokobelanja/app/delivery"
	_repository "Tokobelanja/app/repository"
	_usecase "Tokobelanja/app/usecase"
	"Tokobelanja/config"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config.StartDB()
	db := config.GetDBConnection()

	userRepository := _repository.NewUserRepository(db)
	userUsecase := _usecase.NewUserUsecase(userRepository)

	categoryRepository := _repository.NewCategoryRepository(db)
	categoryUsecase := _usecase.NewCategoryUsecase(categoryRepository)

	productRepository := _repository.NewProductRepository(db)
	productUsecase := _usecase.NewProductUsecase(productRepository, categoryRepository)

	transactionRepository := _repository.NewTransactionRepository(db)
	transactionUsecase := _usecase.NewTransactionUsecase(transactionRepository, userRepository, productRepository, categoryRepository)

	api := router.Group("/")
	_handler.NewUserHandler(api, userUsecase)
	_handler.NewCategoryHandler(api, categoryUsecase)
	_handler.NewProductHandler(api, productUsecase)
	_handler.NewTransactionHandler(api, transactionUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
