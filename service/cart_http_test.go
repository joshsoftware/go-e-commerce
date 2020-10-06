package service

import (
	"encoding/json"
	"errors"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var image_urls = pq.StringArray([]string{"url1", "url2"})

var testCartProduct = []db.CartProduct{
	{
		Id:          1,
		Name:        "abc",
		Quantity:    10,
		Category:    "clothing",
		Price:       2000,
		Description: "abc",
		ImageUrls:   image_urls,
		Discount:    20,
		Tax:         5,
	},
}

type CartHandlerTestSuite struct {
	suite.Suite
	dbMock *db.DBMockStore
}

func (suite *CartHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}
func TestExampleTestSuite(t *testing.T) {
	config.Load()
	suite.Run(t, new(CartHandlerTestSuite))
	suite.Run(t, new(ProductsHandlerTestSuite))
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
				ImageUrls:   image_urls,
				Discount:    20,
				Tax:         5,
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

func (suite *CartHandlerTestSuite) TestAddToCartSuccess() {
	suite.dbMock.On("AddToCart", mock.Anything, 1, 100).Return(1, nil)

	recorder := makeHTTPCall(http.MethodPost,
		"/cart",
		"/cart?productID=100",
		"",
		addToCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"Item added successfully"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestAddToCartProductIDMissingSuccess() {
	recorder := makeHTTPCall(http.MethodPost,
		"/cart",
		"/cart?productID=",
		"",
		addToCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":"product_id missing"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestAddToCartNoRowsSuccess() {
	suite.dbMock.On("AddToCart", mock.Anything, 1, 100).Return(0, nil)

	recorder := makeHTTPCall(http.MethodPost,
		"/cart",
		"/cart?productID=100",
		"",
		addToCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"zero rows affected"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestAddToCartFailure() {
	suite.dbMock.On("AddToCart", mock.Anything, 1, 100).Return(0, errors.New("Error while adding to cart"))

	recorder := makeHTTPCall(http.MethodPost,
		"/cart",
		"/cart?productID=100",
		"",
		addToCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":"could not add item"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestDeleteFromCartSuccess() {
	suite.dbMock.On("DeleteFromCart", mock.Anything, 1, 100).Return(1, nil)

	recorder := makeHTTPCall(http.MethodDelete,
		"/cart",
		"/cart?productID=100",
		"",
		deleteFromCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"Item removed successfully"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestDeleteFromCartProductIDMissingSuccess() {
	recorder := makeHTTPCall(http.MethodDelete,
		"/cart",
		"/cart?productID=",
		"",
		deleteFromCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":"product_id missing"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestDeleteFromCartNoRowsSuccess() {
	suite.dbMock.On("DeleteFromCart", mock.Anything, 1, 100).Return(0, nil)

	recorder := makeHTTPCall(http.MethodDelete,
		"/cart",
		"/cart?productID=100",
		"",
		deleteFromCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"zero rows affected"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestDeleteFromCartFailure() {
	suite.dbMock.On("DeleteFromCart", mock.Anything, 1, 100).Return(0, errors.New("Error while removing from cart"))

	recorder := makeHTTPCall(http.MethodDelete,
		"/cart",
		"/cart?productID=100",
		"",
		deleteFromCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":"could not remove item"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestUpdateIntoCartSuccess() {
	suite.dbMock.On("UpdateIntoCart", mock.Anything, 1, 100, 3).Return(1, nil)

	recorder := makeHTTPCall(http.MethodPut,
		"/cart",
		"/cart?productID=100&quantity=3",
		"",
		updateIntoCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"Quantity updated successfully"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestUpdateIntoCartProductIDMissingSuccess() {
	recorder := makeHTTPCall(http.MethodPut,
		"/cart",
		"/cart?productID=&quantity=3",
		"",
		updateIntoCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":"product_id missing"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestUpdateIntoCartQuantityMissingSuccess() {
	recorder := makeHTTPCall(http.MethodPut,
		"/cart",
		"/cart?productID=100&quantity=",
		"",
		updateIntoCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":"quantity missing"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestUpdateIntoCartNoRowsSuccess() {
	suite.dbMock.On("UpdateIntoCart", mock.Anything, 1, 100, 3).Return(0, nil)

	recorder := makeHTTPCall(http.MethodPut,
		"/cart",
		"/cart?productID=100&quantity=3",
		"",
		updateIntoCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"zero rows affected"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *CartHandlerTestSuite) TestUpdateIntoCartFailure() {
	suite.dbMock.On("UpdateIntoCart", mock.Anything, 1, 100, 3).Return(0, errors.New("Error while updating into cart"))

	recorder := makeHTTPCall(http.MethodPut,
		"/cart",
		"/cart?productID=100&quantity=3",
		"",
		updateIntoCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":"could not update quantity"}`, recorder.Body.String())
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
		[]db.CartProduct{
			db.CartProduct{
				Id:          1,
				Name:        "abc",
				Quantity:    10,
				Category:    "clothing",
				Price:       2000,
				Description: "abc",
				ImageUrls:   image_urls,
				Discount:    20,
				Tax:         5,
			},
		},
		errors.New("Error marshaling cart data"),
	)
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodGet,
		"/cart",
		"/cart",
		"",
		getCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}
