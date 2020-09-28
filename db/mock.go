package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DBMockStore struct {
	mock.Mock
}

//ListUsers mock method
func (m *DBMockStore) ListUsers(ctx context.Context) (users []User, err error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

//GetUser mock method
func (m *DBMockStore) GetUser(ctx context.Context, id int) (user User, err error) {
	args := m.Called(ctx)
	return args.Get(0).(User), args.Error(1)
}

func (m *DBMockStore) DeleteUserByID(ctx context.Context, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DBMockStore) DisableUserByID(ctx context.Context, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}
