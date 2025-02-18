package response

import (
	"strconv"

	"go-products.com/m/internal/product/domain"
)

type ProductResponse struct {
	Sku      string   `json:"sku"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Price    Discount `json:"price"`
}

type Discount struct {
	Original           int     `json:"original"`
	Final              int     `json:"final"`
	DiscountPercentage *string `json:"discount_percentage"`
	Currency           string  `json:"currency"`
}

func FromDomainProducts(products []domain.Product) []ProductResponse {
	productsResponse := make([]ProductResponse, 0)

	for _, product := range products {
		productsResponse = append(productsResponse, fromDomainProduct(product))
	}

	return productsResponse
}

func fromDomainProduct(product domain.Product) ProductResponse {
	discount := product.GetDiscount()
	var discountPercentage *string = nil

	if discount.Percentage != nil {
		discountPercentageValue := strconv.Itoa(int(*discount.Percentage*100)) + "%"
		discountPercentage = &discountPercentageValue
	}

	return ProductResponse{
		Sku:      product.Sku,
		Name:     product.Name,
		Category: product.Category,
		Price: Discount{
			Original:           product.Price,
			Final:              discount.FinalPrice,
			DiscountPercentage: discountPercentage,
			Currency:           product.Currency,
		},
	}
}
