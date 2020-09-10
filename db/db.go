package db

import (
	"context"
)

// Storer - an interface we use to expose methods that do stuff to the underlying database
type Storer interface {
	ListUsers(context.Context) ([]User, error)
	CreateNewUser(context.Context, User) (User, error)
	CheckUserByEmail(context.Context, string) (bool, User, error)
	AuthenticateUser(context.Context, User) (User, error)
	GetUser(context.Context, int) (User, error)
	CreateBlacklistedToken(context.Context, BlacklistedToken) error
	CheckBlacklistedToken(context.Context, string) (bool, int)
}
