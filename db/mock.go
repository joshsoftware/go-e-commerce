package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DBMockStore struct {
	mock.Mock
}

func (m *DBMockStore) ListUsers(ctx context.Context) (users []User, err error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *DBMockStore) AddToCart(ctx context.Context, cartID, productID int) (rowsAffected int64, err error){
	args := m.Called(ctx, cartID, productID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DBMockStore) DeleteFromCart(ctx context.Context, cartID int, productID int) (rowsAffected int64, err error){
	args := m.Called(ctx, cartID, productID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DBMockStore)	UpdateIntoCart(ctx context.Context, cartID int, productID int, quantity int) (rowsAffected int64,err error){
	args := m.Called(ctx, cartID, productID, quantity)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DBMockStore)	AuthenticateUser(ctx context.Context, user User) (validUser User, err error){
	args := m.Called(ctx, user)
	return args.Get(0).(User), args.Error(1)
}

func (m *DBMockStore) CreateBlacklistedToken(ctx context.Context, token BlacklistedToken) (err error){
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *DBMockStore)	CheckBlacklistedToken(ctx context.Context,token string) (status bool,number int){
	args := m.Called(ctx, token)
	return args.Bool(0), args.Int(1)
}

func (m *DBMockStore) GetUser(ctx context.Context,userID int) (user User,err error){
	args := m.Called(ctx, userID)
	return args.Get(0).(User), args.Error(1)
}
