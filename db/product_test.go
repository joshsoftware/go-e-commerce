package db

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductsTestSuite struct {
	suite.Suite
	dbStore Storer
	db      *sqlx.DB
	sqlmock sqlmock.Sqlmock
}

func (suite *ProductsTestSuite) SetupTest() {
	dbStore, dbConn, sqlmock := InitMockDB()
	suite.dbStore = dbStore
	suite.db = dbConn
	suite.sqlmock = sqlmock
}

func (suite *ProductsTestSuite) TearDownTest() {
	suite.db.Close()
}

var urls = []string{"url1", "url2"}

var testProduct = Product{
	Id:           1,
	Name:         "test organization",
	Description:  "test@gmail.com",
	Price:        12,
	Discount:     1,
	Tax:          0.5,
	Quantity:     15,
	CategoryId:   5,
	CategoryName: "2",
	Brand:        "IST",
	Color:        "black",
	Size:         "Medium",
	URLs:         urls,
}

func (suite *ProductsTestSuite) TestCreateNewProductSuccess() {
	product := Product{
		Name:         "test user",
		Description:  "test database",
		Price:        123,
		Discount:     10,
		Tax:          0.5,
		Quantity:     5.0,
		CategoryId:   1,
		CategoryName: "testing",
		Brand:        "testing",
		Color:        "test",
		Size:         "heigh",
		URLs:         []string{"url1", "url2"},
	}

	suite.sqlmock.ExpectBegin()

	suite.sqlmock.ExpectExec("INSERT INTO product").
		WithArgs("test user", "test database", 123, 10, 0.5, 5.0, 1, "testing", "testing", "test", "heigh").
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.sqlmock.ExpectCommit()

	createdProduct, err := suite.dbStore.CreateNewProduct(context.Background(), product)

	assert.Nil(suite.T(), suite.sqlmock.ExpectationsWereMet())
	assert.Equal(suite.T(), createdProduct, product)
	assert.Nil(suite.T(), err)
}

func (suite *ProductsTestSuite) TestCreateNewProductFailure() {
	product := Product{
		Name:         "test user",
		Description:  "test database",
		Price:        123,
		Discount:     10,
		Tax:          0.5,
		Quantity:     5.0,
		CategoryId:   1,
		CategoryName: "testing",
		Brand:        "testing",
		Color:        "test",
		Size:         "heigh",
		URLs:         []string{"url1", "url2"},
	}

	suite.db.Close()
	suite.sqlmock.ExpectBegin()
	suite.sqlmock.ExpectExec("INSERT INTO product").
		WithArgs("test user", "test database", 123, 10, 0.5, 5.0, 1, "testing", "testing", "test", "heigh").
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.sqlmock.ExpectRollback()

	_, err := suite.dbStore.CreateNewProduct(context.Background(), product)

	assert.NotNil(suite.T(), err)
}

func (suite *ProductsTestSuite) TestUpdateProductStockByIdSuccess() {
	product := Product{
		Name:         "test user",
		Description:  "test database",
		Price:        123,
		Discount:     10,
		Tax:          0.5,
		Quantity:     5.0,
		CategoryId:   1,
		CategoryName: "testing",
		Brand:        "testing",
		Color:        "test",
		Size:         "heigh",
		URLs:         []string{"url1", "url2"},
	}

	suite.sqlmock.ExpectBegin()

	UpdatedProduct, err := suite.dbStore.UpdateProductStockById(context.Background(), product, 1)

	assert.Nil(suite.T(), suite.sqlmock.ExpectationsWereMet())
	assert.Equal(suite.T(), UpdatedProduct, product)
	assert.Nil(suite.T(), err)
}

func (suite *ProductsTestSuite) TestUpdateProductStockByIdFailure() {
	product := Product{
		Name:         "test user",
		Description:  "test database",
		Price:        123,
		Discount:     10,
		Tax:          0.5,
		Quantity:     5.0,
		CategoryId:   1,
		CategoryName: "testing",
		Brand:        "testing",
		Color:        "test",
		Size:         "heigh",
		URLs:         []string{"url1", "url2"},
	}

	suite.sqlmock.ExpectBegin()
	updatedProduct, err := suite.dbStore.UpdateProductStockById(context.Background(), product, 1)
	assert.NotEqual(suite.T(), updatedProduct, product)
	assert.NotNil(suite.T(), err)
}
