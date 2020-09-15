package db

import (
	"context"
)

//Storer type interface
type Storer interface {
	ListUsers(ctx context.Context) (userList []User, err error)
	GetUser(ctx context.Context, id int) (user User, err error)
	UpdateUser(ctx context.Context, user User, id int) (err error)
	CreateBlacklistedToken(ctx context.Context, blToken BlacklistedToken) (err error)
	CheckBlacklistedToken(cyx context.Context, token string) (res bool, id int)
	AuthenticateUser(ctx context.Context, u User) (user User, err error)
	//Create(context.Context, User) error

	//Delete(context.Context, string) error
}
