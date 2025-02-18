package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go-products.com/m/internal/product/domain"
	"go-products.com/m/internal/product/infrastructure/handler/response"
	"go-products.com/m/internal/product/use_cases"
	"go-products.com/m/internal/shared/api"
)

func HandleGetProducts(productsRepository domain.ProductRepository) http.HandlerFunc {
	getProductsUseCase := use_cases.NewGetProductsUseCase(productsRepository)

	return func(writer http.ResponseWriter, request *http.Request) {
		filters, err := getProductsFilters(request)
		if err != nil {
			api.InvalidRequest(writer, err.Error())

			return
		}

		ctx := request.Context()
		products, err := getProductsUseCase.Execute(ctx, filters)
		if err != nil {
			api.InternalServerError(writer, err.Error())

			return
		}

		api.Success(writer, response.FromDomainProducts(products))
	}
}

func getProductsFilters(request *http.Request) (domain.ProductsFilters, error) {
	category := api.GetQueryParam(request, "category")
	priceLessThan := api.GetQueryParam(request, "price_less_than")

	priceFilter, err := strconv.Atoi(priceLessThan)
	if priceLessThan != "" && err != nil {
		return domain.ProductsFilters{}, errors.New("price_less_than must be a number")
	}

	categoryFilter := &category
	if category == "" {
		categoryFilter = nil
	}

	limit := 5

	productFilters := domain.ProductsFilters{
		Category:      categoryFilter,
		PriceLessThan: &priceFilter,
		Limit:         &limit,
	}

	if priceLessThan == "" {
		productFilters.PriceLessThan = nil
	}

	return productFilters, nil
}
