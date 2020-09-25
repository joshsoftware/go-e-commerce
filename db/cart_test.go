package db

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CartTestSuite struct {
	suite.Suite
	dbStore Storer
	db      *sqlx.DB
	sqlmock sqlmock.Sqlmock
}

var expectedCart = CartProduct{
	Id:          1,
	Name:        "laptop",
	Quantity:    1,
	Category:    "electronics",
	Price:       202,
	Description: "description",
	ImageUrls:   "urltemp",
}

func (suite *CartTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
	mockedRows = suite.getMockedRows()
}

func (suite *CartTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *CartTestSuite) getMockedRows() (mockedRows *sqlmock.Rows) {
	mockedRows = suite.sqlmock.NewRows([]string{"product_id", "product_name", "quantity", "category_name", "price", "description", "url"}).
		AddRow(1, "laptop", 1, "electronics", 202, "description", "urltemp")
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

}
