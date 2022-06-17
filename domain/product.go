package domain

import (
	"context"
	"time"
)

type Product struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Title      string    `json:"title" gorm:"notNull"`
	Price      int64     `json:"price" gorm:"notNull"`
	Stock      int64     `json:"stock" gorm:"notNull"`
	CategoryID int64     `json:"category_id"  gorm:"notNull"`
	CreatedAt  time.Time `json:"created_at" gorm:"notNull"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"notNull"`
}

type ProductUsecase interface {
	GetProducts(ctx context.Context) (interface{}, error)
	StoreProduct(ctx context.Context, product *Product) (Product, error)
	UpdateProduct(ctx context.Context, product *Product) (Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductRepository interface {
	GetProducts(ctx context.Context) (interface{}, error)
	StoreProduct(ctx context.Context, product *Product) (productId int64, err error)
	GetProductByID(ctx context.Context, id int64) (Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id int64) error
}