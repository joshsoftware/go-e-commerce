package db

import (
	"context"
)

type Storer interface {
	ListProducts(context.Context, string, string) (int, []Product, error)
	FilteredProducts(context.Context, Filter, string, string) (int, []Product, error)
	SearchRecords(context.Context, string, string, string) (int, []Product, error)
	CreateProduct(context.Context, Product) (Product, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int) (Product, error)
	UpdateProductStockById(context.Context, Product, int) (Product, error)
	GetProductByID(context.Context, int) (Product, error)
}
