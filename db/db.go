package db

import (
	"context"
)

type Storer interface {
	ListProducts(context.Context, string, string) ([]Product, error)
	GetProductsByCategoryID(context.Context, int) ([]Product, error)
	CreateNewProduct(context.Context, Product) (Product, error)
	DeleteProductById(context.Context, int) error
	GetProductImagesByID(context.Context, int) ([]ProductImage, error)
	GetProductByID(context.Context, int) (Product, error)
	TotalRecords(context.Context) int
}
