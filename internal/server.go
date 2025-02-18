package internal

import (
	"net/http"

	"go-products.com/m/internal/product/domain"
	"go-products.com/m/internal/product/infrastructure/handler"
	"go-products.com/m/internal/shared/api"
)

func SetupServer(productsRepository domain.ProductRepository) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/products", api.Method(http.MethodGet, handler.HandleGetProducts(productsRepository)))

	return router

}
