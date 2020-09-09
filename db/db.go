package db

import (
	"context"
)

//Storer type interface
type Storer interface {
	ListUsers(context.Context) ([]User, error)
	GetUser(context.Context, int) (User, error)
	UpdateUser(context.Context, User, int) (User, error)
	//Create(context.Context, User) error

	//Delete(context.Context, string) error
}
