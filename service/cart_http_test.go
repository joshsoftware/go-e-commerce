package service

import(
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"errors"
)

type Cart struct {
	cartID int
	productID int
	quantity int
}

var testCart = Cart {
	cartID: 1,
	productID: 100,
	quantity: 1,
}

type CartHandlerTestSuite struct {
	suite.Suite
	dbMock *db.DBMockStore
}

func (suite *CartHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}	
}

func (suite *CartHandlerTestSuite) TestAddToCartSuccess() {
	suite.dbMock.On("AddToCart", mock.Anything, 1, 100).Return(1, nil) 
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPost, 
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
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPost, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPost, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPost, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodDelete, 
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
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodDelete, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodDelete, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodDelete, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPut, 
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
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPut, 
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
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPut, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPut, 
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
	
	recorder := makeHTTPCallWithJWTMiddleware(http.MethodPut, 
		"/cart",
		"/cart?productID=100&quantity=3",
		"",
		updateIntoCartHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":"could not update quantity"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}