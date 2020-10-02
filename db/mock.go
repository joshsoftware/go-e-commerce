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

//deps.store.GetCart(ctx,id)
func (m *DBMockStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	args := m.Called(ctx, user_id)
	return args.Get(0).([]CartProduct), args.Error(1)
}

func (m *DBMockStore) AuthenticateUser(ctx context.Context, u User) (user User, err error) {
	args := m.Called(ctx, u)
	return args.Get(0).(User), args.Error(1)
}

func (m *DBMockStore) CheckBlacklistedToken(ctx context.Context, token string) (bool, int) {
	args := m.Called(ctx, token)
	return args.Get(0).(bool), args.Get(1).(int)
}

func (m *DBMockStore) CreateBlacklistedToken(ctx context.Context, token BlacklistedToken) (err error) {
	args := m.Called(ctx, token)
	return args.Get(0).(error)
}

func (m *DBMockStore) GetUser(ctx context.Context, id int) (user User, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}
