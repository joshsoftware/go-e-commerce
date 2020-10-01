package db

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var (
	now        time.Time
	mockedRows *sqlmock.Rows
)

func InitMockDB() (s Storer, sqlConn *sqlx.DB, sqlmockInstance sqlmock.Sqlmock) {

	// sqlmock.New() gives error : not able to match sql queries ,so adding these parameters : sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual) to sqlmock.New allows complex queries like Join to be matched
	mockDB, sqlmock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.WithField("err:", err).Error("error initializing mock db")
		return
	}

	sqlmockInstance = sqlmock
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	var pgStoreConn pgStore
	pgStoreConn.db = sqlxDB

	return &pgStoreConn, sqlxDB, sqlmockInstance
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CartTestSuite))
}
