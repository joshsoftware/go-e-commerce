package service

import(
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"joshsoftware/go-e-commerce/db"
	"net/http"
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