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

func (suite *CartTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
}

func (suite *CartTestSuite) TestAddToCartSuccess() {
	suite.sqlmock.ExpectExec("INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)").
	WithArgs(1,100).
	WillReturnResult(sqlmock.NewResult(1,1))
	suite.sqlmock.ExpectCommit()

	result, err := suite.dbStore.AddToCart(context.Background(), 1, 100)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, int64(1))
}

func (suite *CartTestSuite) TestAddToCartFailure() {
	suite.db.Close()
	suite.sqlmock.ExpectExec("INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)").
	WithArgs(1,100).
	WillReturnResult(sqlmock.NewResult(1,1))
	suite.sqlmock.ExpectCommit()

	_, err := suite.dbStore.AddToCart(context.Background(), 1, 100)
	assert.NotNil(suite.T(), err)
}

func (suite *CartTestSuite) TestDeleteFromCartSuccess() {
	suite.sqlmock.ExpectExec("DELETE FROM cart WHERE id = $1 AND product_id = $2").
	WithArgs(1,100).
	WillReturnResult(sqlmock.NewResult(1,1))

	result, err := suite.dbStore.DeleteFromCart(context.Background(), 1, 100)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, int64(1))
}

func (suite *CartTestSuite) TestDeleteFromCartFailure() {
	suite.db.Close()
	suite.sqlmock.ExpectExec("DELETE FROM cart WHERE id = $1 AND product_id = $2").
	WithArgs(1,100).
	WillReturnResult(sqlmock.NewResult(1,1))

	_, err := suite.dbStore.DeleteFromCart(context.Background(), 1, 100)
	assert.NotNil(suite.T(), err)
}

func (suite *CartTestSuite) TestUpdateIntoCartSuccess() {
	suite.sqlmock.ExpectExec("UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3").
	WithArgs(3,1,100).
	WillReturnResult(sqlmock.NewResult(1,1))

	result, err := suite.dbStore.UpdateIntoCart(context.Background(), 3, 1, 100)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, int64(1))
}

func (suite *CartTestSuite) TestUpdateIntoCartFailure() {
	suite.db.Close()
	suite.sqlmock.ExpectExec("UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3").
	WithArgs(3,1,100).
	WillReturnResult(sqlmock.NewResult(1,1))

	_, err := suite.dbStore.UpdateIntoCart(context.Background(), 3, 1, 100)
	assert.NotNil(suite.T(), err)
}