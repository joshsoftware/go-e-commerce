package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

// Needs to be implemented using https://github.com/DATA-DOG/go-sqlmock

type ProductTestSuite struct {
	suite.Suite
	dbStore Storer
	db      *sqlx.DB
	sqlmock sqlmock.Sqlmock
}

func (suite *ProductTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
}

func (suite *ProductTestSuite) TearDownTest() {
	suite.db.Close()
}
