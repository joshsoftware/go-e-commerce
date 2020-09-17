package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	CreateUser(context.Context, User) (User, error)
	GetUserByEmail(context.Context, string) (User, error)
}
