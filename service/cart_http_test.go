package service

import (
	"encoding/json"
	"joshsoftware/go-e-commerce/db"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from/     testify - including assertion methods.
var testCartProducts = []db.Product{
	{
		Id:   1,
		Name: "abc",
	},
	{
		Id:   2,
		Name: "pqr",
	},
}

var testcart = []db.Cart{
	{
		Id:        1,
		ProductId: 1,
		Quantity:  1,
	},
	{
		Id:        1,
		ProductId: 2,
		Quantity:  1,
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
	// testGetCart := testCart
	// testGetCart.Id = 1

	suite.dbMock.On("GetCart", mock.Anything, 1).Return(
		testCartProducts,
		nil,
	)

	recorder := makeHTTPCall(http.MethodGet,
		"/user/{user_id:[0-9]+}/cart",
		"/user/1/cart",
		"",
		getCartHandler(Dependencies{Store: suite.dbMock}),
	)

	actual := []db.Product{}
	_ = json.Unmarshal(recorder.Body.Bytes(), &actual)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), testCartProducts, actual)

	suite.dbMock.AssertExpectations(suite.T())

}

// func (suite *CartHandlerTestSuite) TestGetCartDbFailure() {
// 	suite.dbMock.On("GetCart", mock.Anything, 1).Return(
// 		[]db.Product{},
// 		errors.New("Error in fetching data"),
// 	)
// 	recorder := makeHTTPCall(http.MethodGet,
// 		"/users/{cart_id:[0-9]+}/cart",
// 		"/users/1/cart",
// 		"",
// 		getCartHandler(Dependencies{Store: suite.dbMock}),
// 	)
// 	fmt.Println(recorder.Code)
// 	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)

// 	suite.dbMock.AssertExpectations(suite.T())
// }
