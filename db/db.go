package db

import (
	"context"
)

// Storer - an interface we use to expose methods that do stuff to the underlying database
type Storer interface {
	ListUsers(context.Context) ([]User, error)
	CreateNewUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User, int) (User, error)
	CheckUserByEmail(context.Context, string) (bool, User, error)
	AuthenticateUser(context.Context, User) (User, error)
	GetUser(context.Context, int) (User, error)
	CreateBlacklistedToken(context.Context, BlacklistedToken) error
	CheckBlacklistedToken(context.Context, string) (bool, int)
	// TotalRecords(context.Context) (int, error)
	// FilteredRecordsCount(context.Context, Filter) (int, error)
	// FilteredRecords(context.Context, Filter, string, string) ([]Product, error)
	ListProducts(context.Context, string, string) (int, []Product, error)
	FilteredProducts(context.Context, Filter, string, string) (int, []Product, error)
	SearchRecords(context.Context, string, string, string) (int, []Product, error)
	CreateNewProduct(context.Context, Product) (Product, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int) (Product, error)
	UpdateProductStockById(context.Context, Product, int) (Product, error)
	GetProductImagesByID(context.Context, int) ([]ProductImage, error)
	GetProductByID(context.Context, int) (Product, error)
	GetCart(context.Context, int) ([]CartProduct, error)
	AddToCart(context.Context, int, int) (int64, error)
	RemoveFromCart(context.Context, int, int) (int64, error)
	UpdateIntoCart(context.Context, int, int, int) (int64, error)
	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error

	UpdateUserByID(ctx context.Context, user User, id int) (err error)
}
