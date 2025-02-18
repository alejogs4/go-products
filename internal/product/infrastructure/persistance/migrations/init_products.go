package migrations

import (
	"context"
	"database/sql"
	"os"
	"strings"

	"go-products.com/m/internal/product/domain"
)

func InitProducts(ctx context.Context, productsRepository domain.ProductRepository, migrationFilePath string) error {
	fileContent, err := os.ReadFile(migrationFilePath)
	if err != nil {
		return err
	}

	productsCh := ReadJson[domain.CreateProductDTO](fileContent)
	for productDTO := range productsCh {
		if productDTO.Error != nil {
			return err
		}

		if err := productsRepository.CreateProduct(ctx, productDTO.Item); err != nil && !isPrimaryKeyViolation(err) {
			return err
		}
	}

	return nil
}

func isPrimaryKeyViolation(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "UNIQUE constraint") || strings.Contains(errMsg, "PRIMARY KEY")
}

func CreateProductsDatabase(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
    		sku TEXT PRIMARY KEY,
    		name TEXT NOT NULL,
    		category TEXT NOT NULL,
    		price INTEGER NOT NULL
);`)

	return err
}
