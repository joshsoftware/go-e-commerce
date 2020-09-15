package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	CreateUser(context.Context, User) (User, error)
	// CheckUserByEmail(context.Context, string) (bool, User, error)
	GetUserByEmail(context.Context, string) (User, error)
	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error
}
