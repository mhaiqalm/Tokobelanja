package domain

import (
	"context"
	"time"
)

type Transaction struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement:true"`
	ProductID int64 `json:"product_id" gorm:"notNull"`
	UserID int64 `json:"user_id" gorm:"notNull"`
	Quantity int64	`json:"quantity" gorm:"notNull"`
	TotalPrice int64	`json:"total_price" gorm:"notNull"`
	CreatedAt time.Time `json:"created_at" gorm:"notNull"`
	UpdatedAt time.Time `json:"updated_at" gorm:"notNull"`
}

type TransactionUsecase interface {
	GetTransactions(ctx context.Context) (interface{}, error)
	GetMyTransactions(ctx context.Context, id int64) (interface{}, error)
	StoreTransaction(ctx context.Context, transaction *Transaction) (interface{}, error)

}

type TransactionRepository interface {
	GetTransactions(ctx context.Context) (interface{}, error)
	GetMyTransactions(ctx context.Context, userId int64) (interface{}, error)
	StoreTransaction(ctx context.Context, transaction *Transaction) error
}