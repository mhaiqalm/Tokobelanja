package repository

import (
	"context"
	"time"

	"Tokobelanja/domain"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	Conn *gorm.DB
}

func NewTransactionRepository(Conn *gorm.DB) domain.TransactionRepository {
	return &TransactionRepository{Conn}
}

func (t *TransactionRepository) GetTransactions(ctx context.Context) (interface{}, error) {
	type User struct {
		ID        int64     `json:"id"`
		FullName  string    `json:"full_name"`
		Email     string    `json:"email"`
		Balance   int64     `json:"balance"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	type Transaction struct {
		domain.Transaction
		Product domain.Product `json:"Product"`
		User    User           `json:"User"`
	}

	var transactionsProductUser []Transaction
	err := t.Conn.Preload("Product").Preload("User").Find(&transactionsProductUser).Error
	if err != nil {
		return nil, err
	}

	return transactionsProductUser, nil
}

func (t *TransactionRepository) GetMyTransactions(ctx context.Context, userId int64) (interface{}, error) {
	type Transaction struct {
		domain.Transaction
		Product domain.Product `json:"Product"`
	}

	var transactionsProduct []Transaction
	err := t.Conn.Preload("Product").Where("user_id = ?", userId).Find(&transactionsProduct).Error
	if err != nil {
		return nil, err
	}

	return transactionsProduct, nil
}

func (t *TransactionRepository) StoreTransaction(ctx context.Context, transaction *domain.Transaction) error {
	err := t.Conn.Create(transaction).Error
	if err != nil {
		return err
	}
	return nil
}
