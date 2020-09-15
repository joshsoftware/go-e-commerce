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

// CreateUser - test mock
func (m *DBMockStore) CreateUser(ctx context.Context, u User) (user User, err error) {
	args := m.Called(ctx, u)
	return args.Get(0).(User), args.Error(1)
}

// GetUserByEmail - test mock
func (m *DBMockStore) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	args := m.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}
