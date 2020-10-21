package service

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"time"

	"github.com/stretchr/testify/assert"
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

func (suite *SessionsHandlerTestSuite) TestUserLoginSuccess() {
	user := db.User{}
	user.Email = "mayur@gmail.com"
	user.Password = "sagar"

	suite.dbMock.On("AuthenticateUser", mock.Anything, user).Return(user, nil)
	body :=
		`{
			"email" : "mayur@gmail.com",
			"password": "sagar"
		}`
	recorder := makeHTTPCall(
		http.MethodDelete,
		"/login",
		"/login",
		body,
		userLoginHandler(Dependencies{Store: suite.dbMock}),
	)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *SessionsHandlerTestSuite) TestUserLoginFailure() {
	user := db.User{}
	user.Email = "mayur@gmail.com"
	user.Password = "wrongpassword"
	suite.dbMock.On("AuthenticateUser", mock.Anything, user).Return(user, errors.New("Invalid Credentials"))
	body :=
		`{
			"email" : "mayur@gmail.com",
			"password": "wrongpassword"
		}`
	recorder := makeHTTPCall(
		http.MethodDelete,
		"/login",
		"/login",
		body,
		userLoginHandler(Dependencies{Store: suite.dbMock}),
	)
	assert.Equal(suite.T(), http.StatusUnauthorized, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *SessionsHandlerTestSuite) TestUserLogoutSuccess() {
	userBlackListedToken := db.BlacklistedToken{
		UserID:         1,
		ExpirationDate: time.Unix(int64(1605684869), 0),
		Token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDU2ODQ4NjksImlkIjoxLCJpc0FkbWluIjpmYWxzZX0.-rgCTepUicCXyNLL1KiRudxT6NfowuzO2iC4oLn4reg",
	}
	suite.dbMock.On("CreateBlacklistedToken", mock.Anything, userBlackListedToken).Return(nil)

	recorder := makeHTTPCall(
		http.MethodDelete,
		"/logout",
		"/logout",
		"",
		userLogoutHandler(Dependencies{Store: suite.dbMock}),
	)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *SessionsHandlerTestSuite) TestUserLogoutFailure() {
	userBlackListedToken := db.BlacklistedToken{
		UserID:         1,
		ExpirationDate: time.Unix(int64(1605684869), 0),
		Token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDU2ODQ4NjksImlkIjoxLCJpc0FkbWluIjpmYWxzZX0.-rgCTepUicCXyNLL1KiRudxT6NfowuzO2iC4oLn4reg",
	}
	suite.dbMock.On("CreateBlacklistedToken", mock.Anything, userBlackListedToken).Return(errors.New("error blacklisting a token"))

	recorder := makeHTTPCall(
		http.MethodDelete,
		"/logout",
		"/logout",
		"",
		userLogoutHandler(Dependencies{Store: suite.dbMock}),
	)
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *SessionsHandlerTestSuite) TestGenerateJWTSuccess() {
	token, err := generateJwt(1, false)
	assert.Greater(suite.T(), len(token), 0)
	assert.Equal(suite.T(), err, nil)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *SessionsHandlerTestSuite) TestGetDataFromTokenSuccess() {
	userID, exp, isAdmin, _ := getDataFromToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM2ODg4MDYsImlkIjoyLCJpc0FkbWluIjpmYWxzZX0.AzFhRNESL2iyi9xURqeByVU5UaQof9jScUi3RdXakiA")
	assert.Equal(suite.T(), userID, float64(2))
	assert.Equal(suite.T(), isAdmin, false)
	assert.Equal(suite.T(), exp, int64(1603688806))
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *SessionsHandlerTestSuite) TestGetDataFromTokenFailure() {
	userID, exp, isAdmin, err := getDataFromToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM2ODg4MDYsImlkIjoyLCJpc0FkbWluIjpmYWxzZX0.AzFhRNESL2iyi9xURqeByVU5UaQof9jScUi3RdXaki")
	assert.Equal(suite.T(), userID, float64(0))
	assert.Equal(suite.T(), isAdmin, false)
	assert.Equal(suite.T(), exp, int64(0))
	assert.NotEqual(suite.T(), err, nil)
	suite.dbMock.AssertExpectations(suite.T())
}
