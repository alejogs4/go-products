package persistance

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-products.com/m/internal/product/domain"
)

type ProductsSQLiteRepository struct {
	db *sql.DB
}

var (
	ErrGetProducts = errors.New("error getting products")
	ErrParseRow    = errors.New("error parsing product")
)

func NewProductsSQLiteRepository(db *sql.DB) *ProductsSQLiteRepository {
	return &ProductsSQLiteRepository{db: db}
}

func (r *ProductsSQLiteRepository) GetProducts(ctx context.Context, filters domain.ProductsFilters) ([]domain.Product, error) {
	query, params := getQuery(filters)
	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, ErrGetProducts
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.Sku, &product.Name, &product.Category, &product.Price); err != nil {
			return nil, ErrParseRow
		}

		validatedProduct, err := domain.NewProduct(product.Sku, product.Name, product.Category, product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, *validatedProduct)
	}

	return products, nil
}

func (r *ProductsSQLiteRepository) CreateProduct(ctx context.Context, product domain.CreateProductDTO) error {
	domainProduct, err := domain.NewProduct(product.Sku, product.Name, product.Category, product.Price)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO products (sku, name, category, price) VALUES (?, ?, ?, ?);", domainProduct.Sku, domainProduct.Name, domainProduct.Category, domainProduct.Price)

	return err
}

func getQuery(filters domain.ProductsFilters) (string, []interface{}) {
	params := []interface{}{}

	query := strings.Builder{}
	query.WriteString("SELECT sku, name, category, price FROM products")

	if filters.Category != nil || filters.PriceLessThan != nil {
		query.WriteString(" WHERE")
	}

	if filters.Category != nil {
		query.WriteString(" category = ?")
		params = append(params, *filters.Category)
	}

	if filters.Category != nil && filters.PriceLessThan != nil {
		query.WriteString(" AND")
	}

	if filters.PriceLessThan != nil {
		query.WriteString(" price <= ?")
		params = append(params, *filters.PriceLessThan)
	}

	if filters.Limit != nil {
		query.WriteString(" LIMIT ?")
		params = append(params, *filters.Limit)
	}

	query.WriteString(";")

	return query.String(), params
}
