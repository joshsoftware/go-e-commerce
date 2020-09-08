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

func (m *DBMockStore) GetCart(ctx context.Context, user_id int) (cart []Cart, err error) {
	args := m.Called(ctx)
	return args.Get(0).([]Cart), args.Error(1)
}
