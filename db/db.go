package db

import (
	"context"
)

//Storer - interface to add methods used for db operations
type Storer interface {
	ListUsers(ctx context.Context) (user []User, err error)
	GetUser(ctx context.Context, id int) (user User, err error)
	UpdateUserByID(ctx context.Context, user User, id int) (err error)
}
