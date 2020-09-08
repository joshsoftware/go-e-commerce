package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	GetCart(context.Context, int) ([]Cart, error)
	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error
}
