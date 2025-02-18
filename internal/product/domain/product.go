package domain

import (
	"sort"

	"go-products.com/m/internal/product/domain/errors"
)

type Product struct {
	Sku      string
	Name     string
	Category string
	Price    int
	Currency string
}

const EUR = "EUR"

// Discount represents the discount applied to a product, encoding discount function in a function type allows to add new discounts without modifying the Product struct
type discountFn func() *float64

func NewProduct(sku, name, category string, price int) (*Product, error) {
	product := &Product{
		Sku:      sku,
		Name:     name,
		Category: category,
		Price:    price,
		Currency: EUR,
	}

	if err := product.validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) GetDiscount() Discount {
	discount := p.applyDiscounts([]discountFn{
		p.bootsDiscount,
		p.skuDiscount,
	})

	if discount == nil {
		return Discount{
			FinalPrice: p.Price,
			Percentage: nil,
		}
	}

	finalPrice := int(float64(p.Price) * (1 - *discount))
	return Discount{
		FinalPrice: finalPrice,
		Percentage: discount,
	}
}

func (p *Product) applyDiscounts(discountFns []discountFn) *float64 {
	results := make([]float64, 0)
	for _, discountFn := range discountFns {
		if discount := discountFn(); discount != nil {
			results = append(results, *discount)
		}
	}

	if len(results) == 0 {
		return nil
	}

	sort.Float64s(results)
	discount := results[len(results)-1]

	return &discount
}

func (p *Product) bootsDiscount() *float64 {
	categoriesWithDiscount := map[string]float64{
		"boots": 0.3,
	}

	if discount, ok := categoriesWithDiscount[p.Category]; ok {
		return &discount
	}

	return nil
}

func (p *Product) skuDiscount() *float64 {
	skusWithDiscount := map[string]float64{
		"000003": 0.15,
	}

	if discount, ok := skusWithDiscount[p.Sku]; ok {
		return &discount
	}

	return nil
}

func (p *Product) validate() error {
	err := errors.NewNonEmptyString("sku", p.Sku)
	if err != nil {
		return err
	}

	err = errors.NewNonEmptyString("name", p.Name)
	if err != nil {
		return err
	}

	err = errors.NewNonEmptyString("category", p.Category)
	if err != nil {
		return err
	}

	if err = errors.ValidatePrice(p.Price); err != nil {
		return err
	}

	return nil
}
