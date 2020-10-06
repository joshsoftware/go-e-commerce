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

//SetUserPasswordByID mock method
func (m *DBMockStore) SetUserPasswordByID(ctx context.Context, userPassword string, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

//DeleteUserByID mock method
func (m *DBMockStore) DeleteUserByID(ctx context.Context, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DBMockStore) DisableUserByID(ctx context.Context, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DBMockStore) EnableUserByID(ctx context.Context, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DBMockStore) VerifyUserByID(ctx context.Context, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DBMockStore) GetCart(ctx context.Context, user_id int) (products []CartProduct, err error) {
	args := m.Called(ctx, user_id)
	return args.Get(0).([]CartProduct), args.Error(1)
}

// ListUsers - test mock
func (m *DBMockStore) ListProducts(ctx context.Context, limitStr int, pageStr int) (count int, product []Product, err error) {
	args := m.Called(ctx, limitStr, pageStr)
	return args.Get(0).(int), args.Get(1).([]Product), args.Error(2)
}

func (m *DBMockStore) CreateProduct(ctx context.Context, product Product) (productID int, err error) {
	args := m.Called(ctx, product)
	return args.Get(0).(int), args.Error(1)
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

func (m *DBMockStore) UpdateProductById(ctx context.Context, product Product, id int, imageData bool) (updatedProduct Product, err error) {
	args := m.Called(ctx, product, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *DBMockStore) TotalRecords(ctx context.Context) (count int, err error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func (m *DBMockStore) UpdateUserByID(ctx context.Context, user User, id int) (err error) {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DBMockStore) AddToCart(ctx context.Context, cartID, productID int) (rowsAffected int64, err error) {
	args := m.Called(ctx, cartID, productID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DBMockStore) DeleteFromCart(ctx context.Context, cartID int, productID int) (rowsAffected int64, err error) {
	args := m.Called(ctx, cartID, productID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DBMockStore) UpdateIntoCart(ctx context.Context, quantity int, cartID int, productID int) (rowsAffected int64, err error) {
	args := m.Called(ctx, cartID, productID, quantity)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DBMockStore) AuthenticateUser(ctx context.Context, user User) (validUser User, err error) {
	args := m.Called(ctx, user)
	return args.Get(0).(User), args.Error(1)
}

func (m *DBMockStore) CreateBlacklistedToken(ctx context.Context, token BlacklistedToken) (err error) {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *DBMockStore) CheckBlacklistedToken(ctx context.Context, token string) (status bool, number int) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Int(1)
}
