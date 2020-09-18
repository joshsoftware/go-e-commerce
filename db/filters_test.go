package db

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FilterTestSuite struct {
	suite.Suite
	dbStore Storer
	db      *sqlx.DB
	sqlmock sqlmock.Sqlmock
}

func (suite *FilterTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
}

func (suite *FilterTestSuite) TearDownTest() {
	suite.db.Close()
}

var testFilter = Filter{
	CategoryId:   "1",
	Price:        "asc",
	Brand:        "",
	Size:         "",
	Color:        "",
	CategoryFlag: true,
	PriceFlag:    true,
	BrandFlag:    false,
	SizeFlag:     false,
	ColorFlag:    false,
}

func (suite *FilterTestSuite) TestFilteredProductsSuccess() {

	products := []Product{
		Product{
			Id:           2,
			Name:         "Wrangler",
			Description:  "Men  Slim Fit Jeans",
			Price:        600,
			Discount:     20,
			Tax:          12,
			Quantity:     7,
			CategoryId:   1,
			CategoryName: "Clothes",
			Brand:        "Armani",
			Color:        "Charcoal Black",
			Size:         "Large",
			URLs: []string{
				"url1",
				"url2",
			},
		},
		Product{
			Id:           3,
			Name:         "Dragon Jacket",
			Description:  "Made from the skin of one of the dragons",
			Price:        700,
			Discount:     40,
			Tax:          9,
			Quantity:     5,
			CategoryId:   1,
			CategoryName: "Clothes",
			Brand:        "Veteran",
			Color:        "Black",
			Size:         "Extra Large",
			URLs: []string{
				"url1",
				"url2",
			},
		},
	}

	suite.sqlmock.ExpectBegin()
	count, filteredProducts, err := suite.dbStore.FilteredProducts(context.Background(), testFilter, "1", "1")

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, count)
	assert.Equal(suite.T(), products, filteredProducts)
}

func (suite *FilterTestSuite) TestFilteredProductsFailure() {
	suite.sqlmock.ExpectBegin()
	suite.sqlmock.ExpectRollback()

	count, filteredProducts, err := suite.dbStore.FilteredProducts(context.Background(), testFilter, "1", "1")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
	assert.Equal(suite.T(), "", filteredProducts)
}

var text = "xr"

func (suite *FilterTestSuite) TestSearchProductsCountsSuccess() {

	products := []Product{
		Product{
			Id:           9,
			Name:         "Apple iPhone XR (64GB)",
			Description:  "6.1-inch (15.5 cm) Liquid Retina HD LCD display",
			Price:        50000,
			Discount:     6,
			Tax:          15,
			Quantity:     20,
			CategoryId:   3,
			CategoryName: "Mobile",
			Brand:        "Apple",
			Color:        "Grey",
			Size:         "",
			URLs: []string{
				"url1",
				"url2",
			},
		},
	}

	suite.sqlmock.ExpectBegin()
	suite.sqlmock.ExpectExec(` SELECT COUNT(p.id) from products p INNER JOIN category c ON p.category_id = c.id WHERE LOWER(p.name) LIKE LOWER('%xr%') OR LOWER(p.brand) LIKE LOWER('%xr%') OR LOWER(c.name) LIKE LOWER('%xr%')  ;`).
		WillReturnResult(sqlmock.NewResult(1, 2))

	count, searchProducts, err := suite.dbStore.SearchProductsByText(context.Background(), text, "1", "1")

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
	assert.Equal(suite.T(), products, searchProducts)
}

func (suite *FilterTestSuite) TestSearchProductsFailure() {
	suite.sqlmock.ExpectBegin()
	suite.sqlmock.ExpectRollback()

	count, searchProducts, err := suite.dbStore.SearchProductsByText(context.Background(), text, "1", "1")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
	assert.Equal(suite.T(), "", searchProducts)
}
