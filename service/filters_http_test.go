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

	suite.dbMock.On("FilteredProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, []db.Product{}, nil)

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
	/* suite.dbMock.On("FilteredRecordsCount", mock.Anything, mock.Anything).Return(0,
	errors.New("Error getting count of filtered records")) */
	// When calling FilteredRecords with any args, always return
	// that fakeProducts Array along with nil as error
	suite.dbMock.On("FilteredProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, []db.Product{},
		errors.New("Error getting filtered records or Page not Found"))

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products/filters",
		"/products/filters?limit=5&page=1",
		"",
		getProductByFiltersHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), `{"error":{"message":"Error getting filtered records or Page not Found"}}`, recorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())

}

func (suite *FilterHandlerTestSuite) TestGetProductBySearchSuccess() {

	var urls = []string{"url1", "url2"}
	var color = "Black"
	var size = "Medium"

	suite.dbMock.On("SearchProductsByText", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1,
		[]db.Product{
			db.Product{
				Id:           1,
				Name:         "test organization",
				Description:  "test@gmail.com",
				Price:        12,
				Discount:     1,
				Tax:          0.5,
				Quantity:     15,
				CategoryID:   5,
				CategoryName: "2",
				Brand:        "IST",
				Color:        color,
				Size:         size,
				URLs:         urls,
			},
		},
		nil)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products/search",
		"/products/search?limit=5&page=1&text=test",
		"",
		getProductBySearchHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *FilterHandlerTestSuite) TestGetProductBySearchWhenDBFailure() {

	suite.dbMock.On("SearchProductsByText", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, []db.Product{},
		errors.New("Couldn't find any matching search records or Page out of range"))

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products/search",
		"/products/search?limit=5&page=1&text=Apple",
		"",
		getProductBySearchHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), `{"error":{"message":"Couldn't find any matching search records or Page out of range"}}`, recorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}
