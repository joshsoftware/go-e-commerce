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
	suite.dbMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(nil)
	suite.dbMock.On("CheckUserByEmail", mock.Anything, mock.Anything).Return(false, sql.ErrNoRows)
	body :=
		`{
		"first_name" : "test1",
		"last_name" : "test2",
		"email" : "test@gmail.com",
		"password": "password",
		"mobile_number": "8421987845",
		"country": "India",
		"state": "Maharashtra",
		"city": "Nashik",
		"address": "abc"
	}`
	recorder := makeHTTPCall(http.MethodPost,
		"/register",
		body,
		registerUserHandler(Dependencies{Store: suite.dbMock}),
	)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())

}

func (suite *UsersHandlerTestSuite) TestRegisterUserFailure() {
	suite.dbMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(nil)
	suite.dbMock.On("CheckUserByEmail", mock.Anything, mock.Anything).Return(true, sql.ErrNoRows)
	body :=
		`{
		"first_name" : "test1",
		"last_name" : "test2",
		"email" : "test@gmail.com",
		"password": "password",
		"mobile_number": "8421987845",
		"country": "India",
		"state": "Maharashtra",
		"city": "Nashik",
		"address": "abc"
	}`
	recorder := makeHTTPCall(http.MethodPost,
		"/register",
		body,
		registerUserHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), `{"error":"user already registered"}`, recorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	suite.dbMock.AssertNotCalled(suite.T(), "CreateNewUser", mock.Anything, mock.Anything)
	suite.dbMock.AssertCalled(suite.T(), "CheckUserByEmail", mock.Anything, mock.Anything)

}
