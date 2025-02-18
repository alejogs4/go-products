package domain

import "context"

//go:generate moq -out product_repository_mock.go . ProductRepository
type ProductRepository interface {
	GetProducts(ctx context.Context, filters ProductsFilters) ([]Product, error)
	CreateProduct(ctx context.Context, product CreateProductDTO) error
}

type ProductsFilters struct {
	Category      *string
	PriceLessThan *int
	Limit         *int
}

type CreateProductDTO struct {
	Sku      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}
