package service

import (
	"encoding/json"
	"errors"
	"joshsoftware/go-e-commerce/db"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (suite *UsersHandlerTestSuite) TestListUsersSuccess() {
	fakeUser := db.User{}
	faker.FakeData(&fakeUser)

	// Declare an array of db.User and append the fakeUser onto it for use on the dbMock
	fakeUsers := []db.User{}
	fakeUsers = append(fakeUsers, fakeUser)

	suite.dbMock.On("ListUsers", mock.Anything).Return(fakeUsers, nil)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/users",
		"/users",
		"",
		listUsersHandler(Dependencies{Store: suite.dbMock}),
	)

	var users []db.User
	err := json.Unmarshal(recorder.Body.Bytes(), &users)
	if err != nil {
		log.Fatal("Error converting HTTP body from listUsersHandler into User object in json.Unmarshal")
	}

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.NotNil(suite.T(), users[0].ID)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *UsersHandlerTestSuite) TestListUsersWhenDBFailure() {
	suite.dbMock.On("ListUsers", mock.Anything).Return(
		[]db.User{},
		errors.New("error fetching user records"),
	)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/users",
		"/users",
		"",
		listUsersHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *UsersHandlerTestSuite) TestGetUserSuccess() {

	userID := 1
	suite.dbMock.On("GetUser", mock.Anything, userID).Return(
		db.User{
			ID:        1,
			FirstName: "TestUser",
			LastName:  "TestUser",
			Email:     "TestEmail",
			Mobile:    "TestMobile",
			Address:   "Testaddress",
			Password:  "TestPass",
			Country:   "TestCountry",
			State:     "TestState",
			City:      "TestCity",
		}, nil,
	)

	recorder := makeHTTPCall(http.MethodGet,
		"/user",
		"/user",
		"",
		getUserHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{ "data": {
        "id": 1,
        "first_name": "TestUser",
        "last_name": "TestUser",
        "email": "TestEmail",
		"mobile": "TestMobile",
		"address": "Testaddress",
        "password": "TestPass",
        "country": "TestCountry",
        "state": "TestState",
        "city": "TestCity"
       
    }}`, recorder.Body.String())

	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *UsersHandlerTestSuite) TestUpdateUserSuccess() {
	userID := 1
	user := db.User{
		ID:        1,
		FirstName: "UpdatedUserName",
		LastName:  "TestUser",
		Email:     "TestEmail",
		Mobile:    "TestMobile",
		Password:  "TestPass",
		Country:   "TestCountry",
		State:     "TestState",
		City:      "TestCity",
		Address:   "Testaddress",
	}

	suite.dbMock.On("UpdateUserByID", mock.Anything, user, userID).Return(nil)

	body := ` "id": 1,
	"first_name": "UpdatedUserName",
	"last_name": "TestUser",
	"mobile_number": "TestMobile",
	"password": "TestPass",
	"country": "TestCountry",
	"state": "TestState",
	"city": "TestCity",
	"address": "Testaddress"`

	recorder := makeHTTPCall(http.MethodPatch,
		"/user/update",
		"/user/update",
		body,
		updateUserHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data": {
        "id": 1,
        "first_name": "UpdatedUserName",
        "last_name": "TestUser",
        "email": "TestEmail",
		"mobile": "TestMobile",
		"address": "Testaddress",
        "password": "TestPass",
        "country": "TestCountry",
        "state": "TestState",
		"city": "TestCity"
        
    }}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *UsersHandlerTestSuite) TestUpdateUserDbFailure() {
	userID := 1
	user := db.User{
		ID:        1,
		FirstName: "UpdatedUserName",
		LastName:  "TestUser",
		Email:     "TestEmail",
		Mobile:    "TestMobile",
		Password:  "TestPass",
		Country:   "TestCountry",
		State:     "TestState",
		City:      "TestCity",
		Address:   "Testaddress",
	}
	suite.dbMock.On("UpdateUserByID", mock.Anything, user, userID).Return(errors.New("Error while updating user"))

	body := ` "id": 1,
	"first_name": "UpdatedUser",
	"last_name": "TestUser",
	"mobile": "TestMobile",
	"address": "Testaddress",
	"password": "TestPass",
	"country": "TestCountry",
	"state": "TestState",
	"city": "TestCity"
	`

	recorder := makeHTTPCall(http.MethodPut,
		"/user/update",
		"/user/update",
		body,
		updateUserHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func makeHTTPCall(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
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
