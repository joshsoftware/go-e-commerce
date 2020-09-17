package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"time"
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
	fakeUsers := []db.User{
		{
			ID:         1,
			FirstName:  "sagar",
			LastName:   "sonwane",
			Email:      "sagar@gmail.com",
			Mobile:     "8888998887",
			Address:    "abcsde",
			Password:   "sagar11111",
			Country:    "India",
			State:      "MP",
			City:       "Nepa",
			IsAdmin:    true,
			IsDisabled: false,
			CreatedAt:  time.Now(),
		},
		{
			ID:         2,
			FirstName:  "tejas",
			LastName:   "sonwane",
			Email:      "tejas@gmail.com",
			Mobile:     "8888998887",
			Address:    "abcsde",
			Password:   "tejas11111",
			Country:    "India",
			State:      "MP",
			City:       "Nepa",
			IsAdmin:    false,
			IsDisabled: false,
			CreatedAt:  time.Now(),
		},
	}

	suite.dbMock.On("ListUsers", mock.Anything).Return(fakeUsers, nil)
	recorder := makeHTTPCall(http.MethodGet,
		"/users",
		"/users",
		"",
		listUsersHandler(Dependencies{Store: suite.dbMock}))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	// assert.Equal(suite.T(), `{ "data": 		{
	// 	id:         1,
	// 	first_name:  "sagar",
	// 	last_name:   "sonwane",
	// 	email:      "sagar@gmail.com",
	// 	mobile:     "8888998887",
	// 	address:    "abcsde",
	// 	password:   "sagar11111",
	// 	country:    "India",
	// 	state:      "MP",
	// 	city:       "Nepa",
	// 	isAdmin:    true,
	// 	isDisabled: false,
	// 	created_at:  "12:00",
	// },
	// {
	// 	id:         2,
	// 	first_name:  "tejas",
	// 	last_name:   "sonwane",
	// 	email:      "tejas@gmail.com",
	// 	mobile:     "8888998887",
	// 	address:    "abcsde",
	// 	password:   "tejas11111",
	// 	country:    "India",
	// 	state:      "MP",
	// 	city:       "Nepa",
	// 	isAdmin:    false,
	// 	isDisabled: false,
	// 	created_at:  "12:00",
	// }}`, recorder.Body.String())

	suite.dbMock.AssertExpectations(suite.T())
}
