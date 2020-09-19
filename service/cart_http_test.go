package service

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (suite *CartHandlerTestSuite) TestGetCartDbFailure() {
	suite.dbMock.On("GetCart", mock.Anything, mock.Anything).Return(
		[]db.CartProduct{},
		errors.New("Error in fetching data"),
	)
	// suite.dbMock.On("GetCart", mock.Anything, mock.Anything).Return(
	// 	[]db.CartProduct{},
	// 	errors.New(""),
	// )
	recorder := makeHTTPCall(http.MethodGet,
		"/cart",
		"/cart",
		"",
		getCartHandler(Dependencies{Store: suite.dbMock}),
	)
	fmt.Println(recorder.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)

	suite.dbMock.AssertExpectations(suite.T())
}
