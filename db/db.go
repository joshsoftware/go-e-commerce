package db

import (
	"context"
)

type Storer interface {
	ListProducts(context.Context, string, string) (int, []Product, error)
	FilteredProducts(context.Context, Filter, string, string) (int, []Product, error)
	SearchRecords(context.Context, string, string, string) (int, []Product, error)
	CreateNewProduct(context.Context, Product) (Product, error)
	DeleteProductById(context.Context, int) error
	UpdateProductStockById(context.Context, Product, int) (Product, error)
	GetProductImagesByID(context.Context, int) ([]ProductImage, error)
	GetProductByID(context.Context, int) (Product, error)
}
