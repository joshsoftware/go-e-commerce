package service

import (
	"encoding/json"
	"joshsoftware/go-e-commerce/db"
	"log"
	"net/http"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductsHandlerTestSuite struct {
	suite.Suite
	dbMock *db.DBMockStore
}

func (suite *ProductsHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}

func (suite *ProductsHandlerTestSuite) TestListProductsSucess() {

	fakeProduct := db.Product{}
	faker.FakeData(&fakeProduct)

	fakeProducts := []db.Product{}
	fakeProducts = append(fakeProducts, fakeProduct)

	suite.dbMock.On("ListProducts", mock.Anything).Return(fakeProducts, nil)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products",
		"/products",
		"",
		listProductsHandler(Dependencies{Store: suite.dbMock}),
	)

	var products []db.Product
	err := json.Unmarshal(recorder.Body.Bytes(), &products)
	if err != nil {
		log.Fatal("Error converting HTTP body from listUsersHandler into User object in json.Unmarshal")
	}

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.NotNil(suite.T(), products[0].Id)
	suite.dbMock.AssertExpectations(suite.T())
}
