package handler

import (
	"context"
	"embed"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"go-products.com/m/internal/product/domain"
	"go-products.com/m/internal/product/infrastructure/persistance"
	"go-products.com/m/internal/product/infrastructure/persistance/migrations"
	sharedDatabaseUtils "go-products.com/m/internal/shared/database"
)

//go:embed testdata/*.json
var productsContent embed.FS

//go:embed testdata/integration_test/*.json
var productsContentIntegration embed.FS

func TestHandleGetProducts(t *testing.T) {
	assertions := require.New(t)

	tests := []struct {
		name               string
		productsRepository domain.ProductRepository
		expectedStatusCode int
		expectedResponse   string
		priceLessThan      *string
	}{
		{
			name: "Get products successfully returns a 200",
			productsRepository: &domain.ProductRepositoryMock{
				GetProductsFunc: func(ctx context.Context, filters domain.ProductsFilters) ([]domain.Product, error) {
					return []domain.Product{
						{
							Sku:      "0001",
							Name:     "Product 1",
							Category: "sandals",
							Price:    100,
							Currency: domain.EUR,
						},
						{
							Sku:      "0002",
							Name:     "Product 2",
							Category: "boots",
							Price:    100,
							Currency: domain.EUR,
						},
						{
							Sku:      "000003",
							Name:     "Product 3",
							Category: "sandals",
							Price:    100,
							Currency: domain.EUR,
						},
					}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "testdata/successful_response.json",
		},
		{
			name: "Get products with repository error returns a 500",
			productsRepository: &domain.ProductRepositoryMock{
				GetProductsFunc: func(ctx context.Context, filters domain.ProductsFilters) ([]domain.Product, error) {
					return nil, persistance.ErrGetProducts
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "testdata/error_response.json",
		},
		{
			name:               "Get products with bad price returns a 400",
			productsRepository: &domain.ProductRepositoryMock{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "testdata/error_invalid_response.json",
			priceLessThan:      ptr("not_a_number"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			path := "/api/v1/products"
			if tt.priceLessThan != nil {
				path += "?price_less_than=" + *tt.priceLessThan
			}

			request, err := http.NewRequest(http.MethodGet, path, nil)
			assertions.NoError(err)

			handler := HandleGetProducts(tt.productsRepository)

			handler(recorder, request)

			expectedResponse, err := productsContent.ReadFile(tt.expectedResponse)
			assertions.NoError(err)

			assertions.Equal(tt.expectedStatusCode, recorder.Code)
			assertions.JSONEq(string(expectedResponse), recorder.Body.String())
		})
	}
}

func TestIntegration_HandleGetProducts(t *testing.T) {
	assertions := require.New(t)

	testCases := []struct {
		name          string
		assertions    func(writer *httptest.ResponseRecorder)
		category      *string
		priceLessThan *string
	}{
		{
			name: "Get products successfully",
			assertions: func(writer *httptest.ResponseRecorder) {
				productsContent, err := productsContentIntegration.ReadFile("testdata/integration_test/successful_all_response.json")
				assertions.NoError(err)

				assertions.JSONEq(string(productsContent), writer.Body.String())
				assertions.Equal(http.StatusOK, writer.Code)
			},
		},
		{
			name: "Get products with category filter",
			assertions: func(writer *httptest.ResponseRecorder) {
				productsContent, err := productsContentIntegration.ReadFile("testdata/integration_test/successful_category_sandals_response.json")
				assertions.NoError(err)

				assertions.JSONEq(string(productsContent), writer.Body.String())
				assertions.Equal(http.StatusOK, writer.Code)
			},
			category: ptr("sandals"),
		},
		{
			name: "Get products with price less than filter",
			assertions: func(writer *httptest.ResponseRecorder) {
				productsContent, err := productsContentIntegration.ReadFile("testdata/integration_test/successful_price_less_than_75000_response.json")
				assertions.NoError(err)

				assertions.JSONEq(string(productsContent), writer.Body.String())
				assertions.Equal(http.StatusOK, writer.Code)
			},
			priceLessThan: ptr("75000"),
		},
		{
			name: "Get products with category and price less than filter",
			assertions: func(writer *httptest.ResponseRecorder) {
				productsContent, err := productsContentIntegration.ReadFile("testdata/integration_test/successful_category_boots_price_less_than_75000_response.json")
				assertions.NoError(err)

				assertions.JSONEq(string(productsContent), writer.Body.String())
				assertions.Equal(http.StatusOK, writer.Code)
			},
			category:      ptr("boots"),
			priceLessThan: ptr("75000"),
		},
	}

	database, err := sharedDatabaseUtils.GenerateDatabaseConnection(sharedDatabaseUtils.DatabaseConnection{
		DatabaseName: "file::memory:?cache=shared",
	}, migrations.CreateProductsDatabase)
	assertions.NoError(err)

	repository := persistance.NewProductsSQLiteRepository(database)
	err = migrations.InitProducts(context.Background(), repository, path.Join(".", "testdata", "products.json"))
	assertions.NoError(err)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/products"
			if tt.category != nil {
				path += "?category=" + *tt.category
			}

			if tt.priceLessThan != nil {
				if tt.category != nil {
					path += "&price_less_than=" + *tt.priceLessThan
				} else {
					path += "?price_less_than=" + *tt.priceLessThan
				}
			}

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, path, nil)
			assertions.NoError(err)

			handler := HandleGetProducts(repository)

			handler(recorder, request)

			tt.assertions(recorder)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
