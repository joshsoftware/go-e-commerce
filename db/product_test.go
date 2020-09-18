package db

import (
	"context"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (suite *ProductsTestSuite) TestCreateProductSuccess() {
	product := Product{
		Name:         "test user",
		Description:  "test database",
		Price:        123.0,
		Discount:     10.0,
		Tax:          0.5,
		Quantity:     5,
		CategoryId:   1,
		CategoryName: "testing",
		Brand:        "testing",
		Color:        "testing",
		Size:         "testing",
		//URLs:         []string{"url1", "url2"},
	}

	/* suite.sqlmock.ExpectBegin()

	suite.sqlmock.ExpectExec("INSERT INTO products").
		WithArgs("test user", "test database", 123.0, 10.0, 0.5, 5, 1, "testing", "testing", "testing", "testing").
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.sqlmock.ExpectCommit() */

	createdProduct, err := suite.dbStore.CreateProduct(context.Background(), product)
	fmt.Println(createdProduct, err)
	assert.Nil(suite.T(), suite.sqlmock.ExpectationsWereMet())
	assert.Equal(suite.T(), createdProduct, product)
	assert.Nil(suite.T(), err)
}

func (suite *ProductsTestSuite) TestCreateProductFailure() {
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

	_, err := suite.dbStore.CreateProduct(context.Background(), product)

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
		URLs:         pq.StringArray{"url1", "url2"},
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

func (suite *ProductsTestSuite) TestDeleteProductByIdSuccess() {
	suite.sqlmock.ExpectExec("DELETE").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.dbStore.DeleteProductById(context.Background(), testProduct.Id)

	assert.Nil(suite.T(), err)
}

/* func (suite *ProductsTestSuite) TestUpdateProductByIdSuccess() {
	product := Product{
		Name:        "",
		Description: "test database",
		Price:       123,
		Discount:    10,
		Tax:         0.5,
		Quantity:    5,
		CategoryId:  1,
		Brand:       "testing",
		Color:       "test",
		Size:        "heigh",
	}

	//suite.sqlmock.ExpectBegin()

	UpdatedProduct, err := suite.dbStore.UpdateProductById(context.Background(), product, 1)
	//suite.sqlmock.ExpectCommit()
	fmt.Println("Update Product--->", UpdatedProduct)
	assert.Nil(suite.T(), suite.sqlmock.ExpectationsWereMet())
	assert.Equal(suite.T(), UpdatedProduct, product)
	assert.Nil(suite.T(), err)
} */

/* func (suite *ProductsTestSuite) TestUpdateProductByIdSuccess() {
	suite.sqlmock.ExpectExec("UPDATE products").
		WithArgs("test organization", "test@gmail.com", 100.0, 1.0, 2.0, 5, 1, "test", "test", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.sqlmock.ExpectQuery("SELECT").
		WillReturnRows(mockedRows)

	org, err := suite.dbStore.UpdateProductById(context.Background(), testProduct, testProduct.Id)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), testProduct, org)
} */

func (suite *ProductsTestSuite) TestUpdateProductByIdFailure() {
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
	updatedProduct, err := suite.dbStore.UpdateProductById(context.Background(), product, 1)
	assert.NotEqual(suite.T(), updatedProduct, product)
	assert.NotNil(suite.T(), err)
}

/* func (suite *ProductsTestSuite) TestListProductsSuccess() {
	suite.sqlmock.ExpectQuery(getProductQuery).
		WillReturnRows(mockedRows)

	_, org, err := suite.dbStore.ListProducts(context.Background(), "1", "1")

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []Product{testProduct}, org)
} */
