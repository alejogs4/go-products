package errors

import "errors"

var InvalidPrice = errors.New("price must be greater than 0")

func ValidatePrice(price int) error {
	if price <= 0 {
		return InvalidPrice
	}

	return nil
}
