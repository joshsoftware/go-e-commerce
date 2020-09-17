package db

import (
	"context"
)

type Storer interface {
	ListUsers(ctx context.Context) (users []User, err error)
	DeleteUserByID(ctx context.Context, id int) (err error)
}
