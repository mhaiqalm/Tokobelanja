package usecase

import (
	"context"
	"time"

	"Tokobelanja/domain"
)

type ProductUsecase struct {
	productRepository  domain.ProductRepository
	categoryRepository domain.CategoryRepository
}

func NewProductUsecase(p domain.ProductRepository, c domain.CategoryRepository) domain.ProductUsecase {
	return &ProductUsecase{productRepository: p, categoryRepository: c}
}

func (c *ProductUsecase) GetProducts(ctx context.Context) (interface{}, error) {
	products, err := c.productRepository.GetProducts(ctx)
	if err != nil {
		return []domain.Product{}, domain.ErrInternalServerError
	}
	return products, nil
}

func (c *ProductUsecase) StoreProduct(ctx context.Context, product *domain.Product) (domain.Product, error) {
	_, err := c.categoryRepository.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return domain.Product{}, domain.ErrNotFound
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	productId, err := c.productRepository.StoreProduct(ctx, product)
	if err != nil {
		return domain.Product{}, domain.ErrInternalServerError
	}
	product.ID = productId
	return *product, nil
}

func (c *ProductUsecase) UpdateProduct(ctx context.Context, product *domain.Product) (domain.Product, error) {
	_, err := c.productRepository.GetProductByID(ctx, product.ID)
	if err != nil {
		return domain.Product{}, domain.ErrNotFound
	}
	_, err = c.categoryRepository.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return domain.Product{}, domain.ErrNotFound
	}
	product.UpdatedAt = time.Now()
	err = c.productRepository.UpdateProduct(ctx, product)
	if err != nil {
		return domain.Product{}, domain.ErrInternalServerError
	}
	return *product, nil
}

func (c *ProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	_, err := c.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}
	err = c.productRepository.DeleteProduct(ctx, id)
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}
