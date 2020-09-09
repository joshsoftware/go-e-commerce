package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

//DBMockStore contains mock of db
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

//UpdateUser mock method
func (m *DBMockStore) UpdateUser(ctx context.Context, user User, id int) (updatedUser User, err error) {
	args := m.Called(ctx)
	return args.Get(0).(User), args.Error(1)
}
