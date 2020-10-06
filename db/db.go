package db

import (
	"context"
	"mime/multipart"
)

type Storer interface {
	ListProducts(context.Context, int, int) (int, []Product, error)
	FilteredProducts(context.Context, Filter, string, string) (int, []Product, error)
	SearchProductsByText(context.Context, string, string, string) (int, []Product, error)
	CreateProduct(context.Context, Product, []*multipart.FileHeader) (Product, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int, []*multipart.FileHeader) (Product, error, int)
	UpdateProductStockById(context.Context, Product, int) (Product, error)
	GetProductByID(context.Context, int) (Product, error)
}
