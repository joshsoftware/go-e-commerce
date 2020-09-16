package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DBMockStore struct {
	mock.Mock
}

// ListUsers - test mock
func (m *DBMockStore) ListProducts(ctx context.Context, limit string, page string) (count int, product []Product, err error) {
	args := m.Called(ctx, limit, page)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

func (m *DBMockStore) CreateNewProduct(ctx context.Context, p Product) (product Product, err error) {
	args := m.Called(ctx, p)
	return args.Get(0).(Product), args.Error(1)
}

func (m *DBMockStore) FilteredProducts(ctx context.Context, f Filter, limit string, page string) (count int, product []Product, err error) {
	args := m.Called(ctx, f, limit, page)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

func (m *DBMockStore) SearchRecords(ctx context.Context, text string, limit string, page string) (count int, product []Product, err error) {
	args := m.Called(ctx, text, limit, page)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

// DeleteCoreValue - Deletes the core value of the organization
func (m *DBMockStore) DeleteProductById(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DBMockStore) UpdateProductStockById(ctx context.Context, p Product, id int) (product Product, err error) {
	args := m.Called(ctx, p, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *DBMockStore) GetProductImagesByID(ctx context.Context, id int) (productimage []ProductImage, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]ProductImage), args.Error(1)
}

func (m *DBMockStore) GetProductByID(ctx context.Context, id int) (product Product, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Product), args.Error(1)
}

/* func (m *DBMockStore) FilteredRecords(ctx context.Context, filter Filter, limit string, page string) (product []Product, err error) {
	args := m.Called(ctx, filter, limit, page)
	return args.Get(0).([]Product), args.Error(1)
}

func (m *DBMockStore) FilteredRecordsCount(ctx context.Context, filter Filter) (count int, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int), args.Error(1)
}

func (m *DBMockStore) TotalRecords(ctx context.Context) (count int, err error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
} */
