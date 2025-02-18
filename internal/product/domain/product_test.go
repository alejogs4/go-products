package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProduct(t *testing.T) {
	assertions := require.New(t)

	type args struct {
		sku      string
		name     string
		category string
		price    int
	}
	tests := []struct {
		name    string
		args    args
		want    *Product
		wantErr bool
	}{
		{
			name: "Create product successfully",
			args: args{
				sku:      "0001",
				name:     "Product 1",
				category: "sandals",
				price:    100,
			},
			want: &Product{
				Sku:      "0001",
				Name:     "Product 1",
				Category: "sandals",
				Price:    100,
				Currency: EUR,
			},
			wantErr: false,
		},
		{
			name: "Create product with empty sku returns error",
			args: args{
				sku:      "",
				name:     "Product 1",
				category: "sandals",
				price:    100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create product with empty name returns error",
			args: args{
				sku:      "001",
				name:     "",
				category: "sandals",
				price:    100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create product with empty category returns error",
			args: args{
				sku:      "001",
				name:     "name",
				category: "",
				price:    100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create product with empty negative price",
			args: args{
				sku:      "001",
				name:     "name",
				category: "sandals",
				price:    -1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProduct(tt.args.sku, tt.args.name, tt.args.category, tt.args.price)
			assertions.Equal(err != nil, tt.wantErr)
			assertions.Equal(tt.want, got)
		})
	}
}

func TestProduct_GetDiscount(t *testing.T) {
	assertions := require.New(t)

	type fields struct {
		Sku      string
		Name     string
		Category string
		Price    int
		Currency string
	}
	tests := []struct {
		name   string
		fields fields
		want   Discount
	}{
		{
			name: "Get discount for boots",
			fields: fields{
				Sku:      "0001",
				Name:     "Product 1",
				Category: "boots",
				Price:    100,
				Currency: EUR,
			},
			want: Discount{
				FinalPrice: 70,
				Percentage: ptr(0.3),
			},
		},
		{
			name: "Get discount for sku",
			fields: fields{
				Sku:      "000003",
				Name:     "Product 3",
				Category: "sandals",
				Price:    100,
				Currency: EUR,
			},
			want: Discount{
				FinalPrice: 85,
				Percentage: ptr(0.15),
			},
		},
		{
			name: "Get biggest discount when both apply",
			fields: fields{
				Sku:      "000003",
				Name:     "Product 3",
				Category: "boots",
				Price:    100,
				Currency: EUR,
			},
			want: Discount{
				FinalPrice: 70,
				Percentage: ptr(0.3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{
				Sku:      tt.fields.Sku,
				Name:     tt.fields.Name,
				Category: tt.fields.Category,
				Price:    tt.fields.Price,
				Currency: tt.fields.Currency,
			}

			assertions.Equal(tt.want, p.GetDiscount())
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
