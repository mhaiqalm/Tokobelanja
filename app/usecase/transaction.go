package usecase

import (
	"context"
	"fmt"
	"time"

	"Tokobelanja/domain"
)

type TransactionUsecase struct {
	transactionRepository domain.TransactionRepository
	userRepository        domain.UserRepository
	productRepository     domain.ProductRepository
	categoryRepository    domain.CategoryRepository
}

func NewTransactionUsecase(t domain.TransactionRepository, u domain.UserRepository, p domain.ProductRepository, c domain.CategoryRepository) domain.TransactionUsecase {
	return &TransactionUsecase{transactionRepository: t, userRepository: u, productRepository: p, categoryRepository: c}
}

func (t *TransactionUsecase) GetTransactions(ctx context.Context) (interface{}, error) {
	transactions, err := t.transactionRepository.GetTransactions(ctx)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	return transactions, nil
}

func (t *TransactionUsecase) GetMyTransactions(ctx context.Context, id int64) (interface{}, error) {
	transactions, err := t.transactionRepository.GetMyTransactions(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	return transactions, nil
}

func (t *TransactionUsecase) StoreTransaction(ctx context.Context, transaction *domain.Transaction) (interface{}, error) {
	product, err := t.productRepository.GetProductByID(ctx, transaction.ProductID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if transaction.Quantity > product.Stock {
		return nil, domain.ErrStockNotEnough
	}

	user, err := t.userRepository.GetUserByID(ctx, transaction.UserID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if product.Price > user.Balance {
		return nil, domain.ErrBalanceNotEnough
	}

	product.Stock = product.Stock - transaction.Quantity
	product.UpdatedAt = time.Now()
	fmt.Println("cekk", product.Stock)
	err = t.productRepository.UpdateProduct(ctx, &product)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	user.Balance = user.Balance - product.Price
	user.UpdatedAt = time.Now()
	err = t.userRepository.UpdateUser(ctx, &user)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	category, err := t.categoryRepository.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	category.Sold_Product_Amount = category.Sold_Product_Amount + transaction.Quantity
	category.UpdatedAt = time.Now()
	err = t.categoryRepository.UpdateCategory(ctx, &category)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	transaction.TotalPrice = product.Price * transaction.Quantity
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
	err = t.transactionRepository.StoreTransaction(ctx, transaction)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return map[string]interface{}{
		"total_price":   product.Price * transaction.Quantity,
		"quantity":      transaction.Quantity,
		"product_title": product.Title,
	}, nil
}
