package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DBMockStore struct {
	mock.Mock
}

// ListUsers - test mock
func (m *DBMockStore) ListProducts(ctx context.Context, limitStr string, pageStr string) (count int, product []Product, err error) {
	args := m.Called(ctx, limitStr, pageStr)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

func (m *DBMockStore) CreateProduct(ctx context.Context, product Product) (createdProduct Product, err error) {
	args := m.Called(ctx, product)
	return args.Get(0).(Product), args.Error(1)
}

func (m *DBMockStore) FilteredProducts(ctx context.Context, filter Filter, limitStr string, pageStr string) (count int, product []Product, err error) {
	args := m.Called(ctx, filter, limitStr, pageStr)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

func (m *DBMockStore) SearchProductsByText(ctx context.Context, text string, limitStr string, pageStr string) (count int, product []Product, err error) {
	args := m.Called(ctx, text, limitStr, pageStr)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

func (m *DBMockStore) DeleteProductById(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DBMockStore) UpdateProductStockById(ctx context.Context, product Product, id int) (updatedProduct Product, err error) {
	args := m.Called(ctx, product, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *DBMockStore) GetProductByID(ctx context.Context, id int) (product Product, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Product), args.Error(1)
}
