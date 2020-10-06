package db

import (
	"context"

	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CartTestSuite struct {
	suite.Suite
	dbStore Storer
	db      *sqlx.DB
	sqlmock sqlmock.Sqlmock
}

var image_urls = pq.StringArray([]string{"url1", "url2"})

var expectedCart = CartProduct{
	Id:          1,
	Name:        "laptop",
	Quantity:    1,
	Category:    "electronics",
	Price:       202,
	Description: "description",
	ImageUrls:   image_urls,
	Discount:    20,
	Tax:         5,
}

func (suite *CartTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
	mockedRows = suite.getMockedRows()
}

func (suite *CartTestSuite) TestAddToCartSuccess() {
	suite.sqlmock.ExpectExec("INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)").
		WithArgs(1, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := suite.dbStore.AddToCart(context.Background(), 1, 100)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, int64(1))
}

func (suite *CartTestSuite) TestAddToCartFailure() {
	suite.db.Close()
	suite.sqlmock.ExpectExec("INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)").
		WithArgs(1, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.dbStore.AddToCart(context.Background(), 1, 100)
	assert.NotNil(suite.T(), err)
}

func (suite *CartTestSuite) TestDeleteFromCartSuccess() {
	suite.sqlmock.ExpectExec("DELETE FROM cart WHERE id = $1 AND product_id = $2").
		WithArgs(1, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := suite.dbStore.DeleteFromCart(context.Background(), 1, 100)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, int64(1))
}

func (suite *CartTestSuite) TestDeleteFromCartFailure() {
	suite.db.Close()
	suite.sqlmock.ExpectExec("DELETE FROM cart WHERE id = $1 AND product_id = $2").
		WithArgs(1, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.dbStore.DeleteFromCart(context.Background(), 1, 100)
	assert.NotNil(suite.T(), err)
}

func (suite *CartTestSuite) TestUpdateIntoCartSuccess() {
	suite.sqlmock.ExpectExec("UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3").
		WithArgs(3, 1, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := suite.dbStore.UpdateIntoCart(context.Background(), 3, 1, 100)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, int64(1))
}

func (suite *CartTestSuite) TestUpdateIntoCartFailure() {
	suite.db.Close()
	suite.sqlmock.ExpectExec("UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3").
		WithArgs(3, 1, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.dbStore.UpdateIntoCart(context.Background(), 3, 1, 100)
	assert.NotNil(suite.T(), err)
}

func (suite *CartTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *CartTestSuite) getMockedRows() (mockedRows *sqlmock.Rows) {
	mockedRows = suite.sqlmock.NewRows([]string{"product_id", "product_name", "quantity", "category_name", "price", "description", "image_urls", "discount", "tax"}).
		AddRow(1, "laptop", 1, "electronics", 202, "description", image_urls, 20, 5)
	return
}

func (suite *CartTestSuite) TestGetCartSuccess() {

	suite.sqlmock.ExpectQuery(joinCartProductQuery).
		WithArgs(1).
		WillReturnRows(mockedRows)

	cart, err := suite.dbStore.GetCart(context.Background(), expectedCart.Id)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []CartProduct{expectedCart}, cart)
}

func (suite *CartTestSuite) TestCartFailure() {

	suite.db.Close()
	//Close connection to test failure case

	suite.sqlmock.ExpectQuery(joinCartProductQuery).
		WillReturnRows(mockedRows)

	_, err := suite.dbStore.GetCart(context.Background(), expectedCart.Id)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.New("sql: database is closed"), err)
}
