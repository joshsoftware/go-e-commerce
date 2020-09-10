package db

import (
	"context"
)

//Storer type interface
type Storer interface {
	ListUsers(context.Context) ([]User, error)
	GetUser(context.Context, int) (User, error)
	UpdateUser(context.Context, User, int) (User, error)
	CreateBlacklistedToken(context.Context, BlacklistedToken) error
	CheckBlacklistedToken(context.Context, string) (bool, int)
	AuthenticateUser(context.Context, User) (User, error)
	//Create(context.Context, User) error

	//Delete(context.Context, string) error
}
