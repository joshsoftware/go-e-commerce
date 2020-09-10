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

// CreateNewUser - test mock
func (m *DBMockStore) CreateNewUser(ctx context.Context, u User) (user User, err error) {
	args := m.Called(ctx, u)
	return args.Get(0).(User), args.Error(1)
}

// CheckUserByEmail - test mock
func (m *DBMockStore) CheckUserByEmail(ctx context.Context, email string) (check bool, user User, err error) {
	args := m.Called(ctx, email)
	return args.Get(0).(bool), args.Get(1).(User), args.Get(2).(error)
}
