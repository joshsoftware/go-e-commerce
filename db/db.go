package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	AddToCart(context.Context, int, int) (error)
	RemoveFromCart(context.Context, int, int) (error)
	UpdateIntoCart(context.Context, int, int, int) (error)
	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error
}
