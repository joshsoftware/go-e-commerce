package service

import (
	"encoding/json"
	"errors"
	"joshsoftware/go-e-commerce/db"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var testCartProduct = []db.CartProduct{
	{
		Id:          1,
		Name:        "abc",
		Quantity:    10,
		Category:    "clothing",
		Price:       2000,
		Description: "abc",
		ImageUrls:   "temp",
	},
}

type CartHandlerTestSuite struct {
	suite.Suite

	dbMock *db.DBMockStore
}

func (suite *CartHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}

func (suite *CartHandlerTestSuite) TestGetCartSuccess() {
	suite.dbMock.On("GetCart", mock.Anything, mock.Anything).Return(
		[]db.CartProduct{
			db.CartProduct{
				Id:          1,
				Name:        "abc",
				Quantity:    10,
				Category:    "clothing",
				Price:       2000,
				Description: "abc",
				ImageUrls:   "temp",
			},
		},
		nil,
	)

	recorder := makeHTTPCallWithJWTMiddleware(http.MethodGet,
		"/cart",
		"/cart",
		"",
		getCartHandler(Dependencies{Store: suite.dbMock}),
	)

	actual := []db.CartProduct{}
	_ = json.Unmarshal(recorder.Body.Bytes(), &actual)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), testCartProduct, actual)

	suite.dbMock.AssertExpectations(suite.T())

}

func (suite *CartHandlerTestSuite) TestGetCartDbFailureFetchError() {
	suite.dbMock.On("GetCart", mock.Anything, mock.Anything).Return(
		[]db.CartProduct{},
		errors.New("Error fetching data from database"),
	)

	recorder := makeHTTPCallWithJWTMiddleware(http.MethodGet,
		"/cart",
		"/cart",
		"",
		getCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":"Error fetching data from database"}`, recorder.Body.String())

	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestGetCartDbFailureJSONMarshallError() {
	suite.dbMock.On("GetCart", mock.Anything, mock.Anything).Return(
		[]db.CartProduct{},
		// x.
		errors.New("Error marshaling cart data"),
	)
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodGet,
		"/cart",
		"/cart",
		"",
		getCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	// assert.Equal(suite.T(), `{"error":"Error marshaling cart data"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

// func (suite *CartHandlerTestSuite) TestGetCartDbFailureAuthorizationError() {
// 	suite.dbMock.On("GetCart", mock.Anything, mock.Anything).Return(
// 		[]db.CartProduct{},
// 		errors.New("Error fetching data from token"),
// 	)
// 	recorder := makeHTTPCallWithJWTMiddleware(http.MethodGet,
// 		"/cart",
// 		"/cart",
// 		"",
// 		getCartHandler(Dependencies{Store: suite.dbMock}),
// 	)
// 	fmt.Println(recorder.Code, "Authooo", http.StatusUnauthorized)
// 	assert.Equal(suite.T(), http.StatusUnauthorized, recorder.Code)

// 	suite.dbMock.AssertExpectations(suite.T())
// }
