package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	CreateNewUser(context.Context, User) error
	CheckUserByEmail(context.Context, string) (bool, error)
	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error
}
