package use_cases

import (
	"context"

	"go-products.com/m/internal/product/domain"
)

type GetProductsUseCase struct {
	productRepository domain.ProductRepository
}

func NewGetProductsUseCase(productRepository domain.ProductRepository) GetProductsUseCase {
	return GetProductsUseCase{productRepository: productRepository}
}

func (u GetProductsUseCase) Execute(ctx context.Context, filters domain.ProductsFilters) ([]domain.Product, error) {
	return u.productRepository.GetProducts(ctx, filters)
}
