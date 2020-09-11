package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	AuthenticateUser(context.Context, User) (User, error)
	GetUser(context.Context, int) (User, error)
	CreateBlacklistedToken(context.Context, BlacklistedToken) error
	CheckBlacklistedToken(context.Context, string) (bool, int)
	GetCart(context.Context, int) ([]CartProduct, error)
}
