package service

import (
	"joshsoftware/go-e-commerce/db"

	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type SessionsHandlerTestSuite struct {
	suite.Suite

	dbMock *db.DBMockStore
}

func (suite *SessionsHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}
