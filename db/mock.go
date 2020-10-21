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

// AuthenticateUser function
func (m *DBMockStore) AuthenticateUser(ctx context.Context, u User) (user User, err error) {
	args := m.Called(ctx, u)
	return args.Get(0).(User), args.Error(1)
}

// CreateBlacklistedToken function
func (m *DBMockStore) CreateBlacklistedToken(ctx context.Context, token BlacklistedToken) (err error) {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// GetUser function
func (m *DBMockStore) GetUser(ctx context.Context, userID int) (user User, err error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(User), args.Error(1)
}

// CheckBlacklistedToken function
func (m *DBMockStore) CheckBlacklistedToken(ctx context.Context, token string) (isBlackListed bool, userID int) {
	args := m.Called(ctx, isBlackListed)
	return args.Get(0).(bool), args.Get(1).(int)
}
