package service

import (
	"database/sql"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type UsersHandlerTestSuite struct {
	suite.Suite

	dbMock *db.DBMockStore
}

func (suite *UsersHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(UsersHandlerTestSuite))
}

func makeHTTPCall(method, path, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create a http request using the given parameters
	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	// test recorder created for capturing api responses
	recorder = httptest.NewRecorder()

	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}

func (suite *UsersHandlerTestSuite) TestRegisterUserSuccess() {
	user := db.User{}
	user.Email = "test@gmail.com"
	user.FirstName = "test1"
	user.LastName = "test2"
	user.Mobile = "8421987856"
	user.Address = "abc"
	user.State = "Maharashtra"
	user.City = "Nashik"
	user.Password = "password"
	user.Country = "India"

	suite.dbMock.On("CreateUser", mock.Anything, user).Return(user, nil)
	suite.dbMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(user, sql.ErrNoRows)
	body :=
		`{
			"first_name" : "test1",
			"last_name" : "test2",
			"email" : "test@gmail.com",
			"mobile": "8421987856",
			"country": "India",
			"state": "Maharashtra",
			"city": "Nashik",
			"address": "abc",
			"password": "password"
		}`
	recorder := makeHTTPCall(http.MethodPost,
		"/register",
		body,
		registerUserHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())

}

func (suite *UsersHandlerTestSuite) TestRegisterUserFailure() {
	user := db.User{}
	user.Email = "test@gmail.com"
	user.FirstName = "test1"
	user.LastName = "test2"
	user.Mobile = "8421987856"
	user.Address = "abc"
	user.State = "Maharashtra"
	user.City = "Nashik"
	user.Password = "password"
	user.Country = "India"

	suite.dbMock.On("CreateUser", mock.Anything, user).Return(user, nil)
	suite.dbMock.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil)
	body :=
		`{
		"first_name" : "test1",
		"last_name" : "test2",
		"email" : "test@gmail.com",
		"mobile": "8421987856",
		"country": "India",
		"state": "Maharashtra",
		"city": "Nashik",
		"address": "abc",
		"password": "password"
	}`
	recorder := makeHTTPCall(http.MethodPost,
		"/register",
		body,
		registerUserHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), `{"error":"user already registered"}`, recorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	suite.dbMock.AssertNotCalled(suite.T(), "CreateUser", mock.Anything, user)
	suite.dbMock.AssertCalled(suite.T(), "GetUserByEmail", mock.Anything, user.Email)

}
