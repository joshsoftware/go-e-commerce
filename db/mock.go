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

func (m *DBMockStore) GetCart(ctx context.Context, user_id int) (products []Product, err error) {
	args := m.Called(ctx, user_id)
	return args.Get(0).([]Product), args.Error(1)
}

// ListUsers - test mock
func (m *DBMockStore) ListProducts(ctx context.Context, limit string, page string) (product []Product, err error) {
	args := m.Called(ctx, limit, page)
	return args.Get(0).([]Product), args.Error(1)
}

func (m *DBMockStore) CreateNewProduct(ctx context.Context, p Product) (product Product, err error) {
	args := m.Called(ctx, p)
	return args.Get(0).(Product), args.Error(1)
}

// DeleteCoreValue - Deletes the core value of the organization
func (m *DBMockStore) DeleteProductById(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DBMockStore) FilteredRecords(ctx context.Context, filter Filter, limit string, page string) (product []Product, err error) {
	args := m.Called(ctx, filter, limit, page)
	return args.Get(0).([]Product), args.Error(1)
}

func (m *DBMockStore) FilteredRecordsCount(ctx context.Context, filter Filter) (count int, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int), args.Error(1)
}

func (m *DBMockStore) GetProductByID(ctx context.Context, id int) (product Product, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *DBMockStore) TotalRecords(ctx context.Context) (count int, err error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func (m *DBMockStore) GetProductImagesByID(ctx context.Context, id int) (productimage []ProductImage, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]ProductImage), args.Error(1)
}

func (m *DBMockStore) UpdateProductStockById(ctx context.Context, p Product, id int) (product Product, err error) {
	args := m.Called(ctx, p, id)
	return args.Get(0).(Product), args.Error(1)
}
