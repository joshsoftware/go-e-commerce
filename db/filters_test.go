package db

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FilterTestSuite struct {
	suite.Suite
	dbStore Storer
	db      *sqlx.DB
	sqlmock sqlmock.Sqlmock
}

func (suite *FilterTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
}

func (suite *FilterTestSuite) TearDownTest() {
	suite.db.Close()
}

var testFilter = Filter{
	CategoryId:   "1",
	Price:        "120",
	Brand:        "test",
	Size:         "Larger",
	Color:        "black",
	CategoryFlag: true,
	PriceFlag:    true,
	BrandFlag:    true,
	SizeFlag:     true,
	ColorFlag:    true,
}

func (suite *FilterTestSuite) TestFilteredProductsSuccess() {
	suite.sqlmock.ExpectQuery(getFilterProduct).
		WillReturnRows(mockedRows)

	_, org, err := suite.dbStore.FilteredProducts(context.Background(), testFilter, "1", "1")

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []Filter{testFilter}, org)
}
