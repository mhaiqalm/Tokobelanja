package repository

import (
	"context"

	"Tokobelanja/domain"

	"gorm.io/gorm"
)

type ProductRepository struct {
	Conn *gorm.DB
}

func NewProductRepository(Conn *gorm.DB) domain.ProductRepository {
	return &ProductRepository{Conn}
}

func (p *ProductRepository) GetProducts(ctx context.Context) (interface{}, error) {

	var products []domain.Product
	err := p.Conn.Find(&products).Error
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

func (p *ProductRepository) StoreProduct(ctx context.Context, product *domain.Product) (productId int64, err error) {
	err = p.Conn.Create(product).Error
	if err != nil {
		return
	}
	productId = product.ID
	return
}

func (p *ProductRepository) GetProductByID(ctx context.Context, id int64) (domain.Product, error) {
	var product domain.Product
	err := p.Conn.First(&product, "id=?", id).Error
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	err := p.Conn.Model(product).Updates(product).Error
	return err
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	var product domain.Product
	product.ID = id
	err := p.Conn.Delete(&product).Error
	return err
}
