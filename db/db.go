package db

import (
	"context"
)

// Storer - an interface we use to expose methods that do stuff to the underlying database
type Storer interface {
	// ListUsers(context.Context) ([]User, error)
	CreateNewUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User, int) (User, error)
	CheckUserByEmail(context.Context, string) (bool, User, error)
	AuthenticateUser(context.Context, User) (User, error)
	CreateBlacklistedToken(context.Context, BlacklistedToken) error
	CheckBlacklistedToken(context.Context, string) (bool, int)
	UpdateUserByID(ctx context.Context, user User, id int) (err error)
	ListUsers(ctx context.Context) (users []User, err error)
	GetUser(ctx context.Context, id int) (user User, err error)
	DeleteUserByID(ctx context.Context, id int) (err error)
	DisableUserByID(ctx context.Context, id int) (err error)
	EnableUserByID(ctx context.Context, id int) (err error)
	VerifyUserByID(ctx context.Context, id int) (err error)
	SetUserPasswordByID(ctx context.Context, password string, id int) (err error)

	// product related
	ListProducts(context.Context, int, int) (int, []Product, error)
	FilteredProducts(context.Context, Filter, string, string) (int, []Product, error)
	SearchProductsByText(context.Context, string, string, string) (int, []Product, error)
	CreateProduct(context.Context, Product) (int, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int, bool) (Product, error)
	UpdateProductStockById(context.Context, Product, int) (Product, error)
	GetProductByID(context.Context, int) (Product, error)

	// cart related
	GetCart(context.Context, int) ([]CartProduct, error)
	AddToCart(context.Context, int, int) (int64, error)
	RemoveFromCart(context.Context, int, int) (int64, error)
	UpdateIntoCart(context.Context, int, int, int) (int64, error)
}
