package service

import (
	"errors"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify, including assertion methods.

type FilterHandlerTestSuite struct {
	suite.Suite

	dbMock *db.DBMockStore
}

func (suite *FilterHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}

func TestFilterHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FilterHandlerTestSuite))
}

// function covers FilteredRecordsCount as well as Filteredrecords
func (suite *FilterHandlerTestSuite) TestGetProductByFiltersSuccess() {

	suite.dbMock.On("FilteredRecordsCount", mock.Anything, mock.Anything).Return(1, nil)

	suite.dbMock.On("FilteredRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]db.Product{}, nil)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products/filters",
		"/products/filters?limit=5&page=1&brand=Apple&categoryid=;&brand=Apple&color=Black",
		"",
		getProductByFiltersHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *FilterHandlerTestSuite) TestFilteredRecordsWhenDBFailure() {

	// Count not expected on filter with brand as Apple and price in desc as failure test
	suite.dbMock.On("FilteredRecordsCount", mock.Anything, mock.Anything).Return(0,
		errors.New("Error getting count of filtered records"))
	// When calling FilteredRecords with any args, always return
	// that fakeProducts Array along with nil as error
	suite.dbMock.On("FilteredRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]db.Product{}, nil)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products/filters",
		"/product/filters?limit=5&page=1&brand=Apple&price=desc",
		"",
		getProductByFiltersHandler(Dependencies{Store: suite.dbMock}),
	)

	//assert.Equal(suite.T(), `{"error":"Error getting count of filtered records"}`, recorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	//suite.dbMock.AssertNotCalled(suite.T(), "FilteredRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	//suite.dbMock.AssertCalled(suite.T(), "FilteredRecordsCount", mock.Anything, mock.Anything)
}
