package service

import (
	"bytes"
	"errors"
	"fmt"
	"joshsoftware/go-e-commerce/db"
	"mime/multipart"
	"net/http"
	"testing"

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

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ProductsHandlerTestSuite))
}

func (suite *ProductsHandlerTestSuite) TestGetProductByIdHandlerSuccess() {

	suite.dbMock.On("GetProductByID", mock.Anything, mock.Anything).Return(
		db.Product{
			Id:           1,
			Name:         "test",
			Description:  "test database",
			Price:        123,
			Discount:     10,
			Tax:          0.5,
			Quantity:     5.0,
			CategoryId:   1,
			CategoryName: "testing",
			Brand:        "new brand",
			Color:        "Black",
			Size:         "Larger",
		}, nil,
	)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products/{product_id:[0-9]+}",
		"/products/1",
		"",
		getProductByIdHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"id":1,"product_title":"test","description":"test database","product_price":123,"discount":10,"tax":0.5,"stock":5,"category_id":1,"category":"testing","brand":"new brand","color":"Black","size":"Larger","image_urls":null}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())

}

func (suite *ProductsHandlerTestSuite) TestGetProductByIdWhenDBFailure() {

	suite.dbMock.On("GetProductByID", mock.Anything, mock.Anything).Return(
		db.Product{}, errors.New("Error in fetching data"),
	)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/product/{product_id:[0-9]+}",
		"/product/1",
		"",
		getProductByIdHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":{"message":"Error feching data, no row found."}}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *ProductsHandlerTestSuite) TestListProductsSuccess() {

	//suite.dbMock.On("TotalRecords", mock.Anything).Return(2, nil)
	suite.dbMock.On("ListProducts", mock.Anything, mock.Anything, mock.Anything).Return(1,
		[]db.Product{
			db.Product{
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
				Color:        "Black",
				Size:         "Larger",
			},
		},
		nil,
	)

	recorder := makeHTTPCall(
		http.MethodGet,
		"/products",
		"/products?limit=1&page=1",
		"",
		listProductsHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"products":[{"id":1,"product_title":"test organization","description":"test@gmail.com","product_price":12,"discount":1,"tax":0.5,"stock":15,"category_id":5,"category":"2","brand":"IST","color":"Black","size":"Larger","image_urls":null}],"total_pages":1}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *ProductsHandlerTestSuite) TestListProductsDBFailure() {

	suite.dbMock.On("ListProducts", mock.Anything, mock.Anything, mock.Anything).Return(0,
		[]db.Product{},
		errors.New("error fetching Products records"),
	)

	recorder := makeHTTPCall(http.MethodGet,
		"/products",
		"/products?limit=1&page=1",
		"",
		listProductsHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":{"message":"Couldn't find any Products records or Page out of range"}}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

var urls = []string{"url1", "url2"}
var color = "Black"
var size = "Medium"

var testProduct = db.Product{
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
	Color:        "Black",
	Size:         "Medium",
	URLs:         urls,
}

func (suite *ProductsHandlerTestSuite) TestCreateProductSuccess() {

	//suite.dbMock.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(testProduct, nil)
	//url := "http://localhost:33001/createProduct"
	//method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("product_title", "test organization")
	_ = writer.WriteField("description", "test@gmail.com")
	_ = writer.WriteField("product_price", "12")
	_ = writer.WriteField("discount", "1")
	_ = writer.WriteField("tax", "0.5")
	_ = writer.WriteField("stock", "15")
	_ = writer.WriteField("category_id", "5")
	_ = writer.WriteField("category", "2")
	_ = writer.WriteField("color", "Black")
	_ = writer.WriteField("size", "Medium")
	_ = writer.WriteField("brand", "IST")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	body := payload.String()
	fmt.Println("body: ", body)

	suite.dbMock.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(
		db.Product{
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
			Color:        "Black",
			Size:         "Medium",
			URLs:         urls,
		},
		nil,
	)

	recorder := makeHTTPCallWithHeader(
		http.MethodPost,
		"/createProduct",
		"/createProduct",
		writer,
		payload,
		createProductHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":{"id":1,"product_title":"test organization","description":"test@gmail.com","product_price":12,"discount":1,"tax":0.5,"stock":15,"category_id":5,"category":"2","brand":"IST","color":"Black","size":"Medium","image_urls":["url1","url2"]}}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *ProductsHandlerTestSuite) TestCreateProductFailure() {

	/* suite.dbMock.On("CreateProduct", mock.Anything, mock.Anything).Return(db.Product{}, sql.ErrNoRows)
		body := `{
			"product_title": "Lancer new",
	        "description": "Mens Running Shoes",
	        "product_price": 150,
			"discount": 15,
			"Tax":          0.5,
	        "stock": 10,
	        "category_id": 6,
	        "category": "Sports",
	        "image_urls": [
	            "Lancer1.jpg",
	            "Lancer2.jpg",
	            "Lancer3.jpg"
	        ]
		}`

		recorder := makeHTTPCall(
			http.MethodPost,
			"/createProduct",
			"/createProduct",
			body,
			createProductHandler(Dependencies{Store: suite.dbMock}),
		) */

	//url := "http://localhost:33001/createProduct"
	//method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("product_title", "test organization6")
	_ = writer.WriteField("description", "test@gmail.com")
	_ = writer.WriteField("product_price", "12")
	_ = writer.WriteField("discount", "1")
	_ = writer.WriteField("tax", "0.5")
	_ = writer.WriteField("stock", "15")
	_ = writer.WriteField("category_id", "5")
	_ = writer.WriteField("category", "2")
	_ = writer.WriteField("color", "Black")
	_ = writer.WriteField("size", "Medium")
	_ = writer.WriteField("brand", "IST")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	body := payload.String()
	fmt.Println("body: ", body)

	suite.dbMock.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(db.Product{},
		errors.New("Error inserting the product, product already exists"))

	recorder := makeHTTPCallWithHeader(
		http.MethodPost,
		"/createProduct",
		"/createProduct",
		writer,
		payload,
		createProductHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), `{"error":{"message":"Error inserting the product, product already exists"}}`, recorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *ProductsHandlerTestSuite) TestCreateProductValidationFailure() {

	/* body := `{
			"product_title": "",
	        "description": "",
	        "product_price": 150,
			"discount": 15,
			"Tax":          0.5,
	        "stock": 10,
	        "category_id": 6,
	        "category": "Sports",
	        "image_url": [
	            "Lancer1.jpg",
	            "Lancer2.jpg",
	            "Lancer3.jpg"
	        ]
		}`

		recorder := makeHTTPCall(http.MethodPost,
			"/createProduct",
			"/createProduct",
			body,
			createProductHandler(Dependencies{Store: suite.dbMock}),
		)

	url := "http://localhost:33001/createProduct"
	method := "POST" */

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("product_title", "test organization6")
	_ = writer.WriteField("description", "")
	_ = writer.WriteField("product_price", "12")
	_ = writer.WriteField("discount", "1")
	_ = writer.WriteField("tax", "0.5")
	_ = writer.WriteField("stock", "15")
	_ = writer.WriteField("category_id", "5")
	_ = writer.WriteField("category", "2")
	_ = writer.WriteField("color", "Black")
	_ = writer.WriteField("size", "Medium")
	_ = writer.WriteField("brand", "IST")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	body := payload.String()
	fmt.Println("body: ", body)

	suite.dbMock.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(db.Product{},
		errors.New("Error inserting the product, product already exists"))

	recorder := makeHTTPCallWithHeader(
		http.MethodPost,
		"/createProduct",
		"/createProduct",
		writer,
		payload,
		createProductHandler(Dependencies{Store: suite.dbMock}),
	)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":{"code":"Invalid_data","message":"Please Provide valid Product data","fields":{"product_description":"Can't be blank "}}}`, recorder.Body.String())
}

func (suite *ProductsHandlerTestSuite) TestDeleteProductByIdSuccess() {

	suite.dbMock.On("DeleteProductById", mock.Anything, 1).Return(
		nil,
	)

	recorder := makeHTTPCall(http.MethodDelete,
		"/product/{product_id:[0-9]+}",
		"/product/1",
		"",
		deleteProductByIdHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":"Product Deleted Successfully!"}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *ProductsHandlerTestSuite) TestDeleteProductByIdDbFailure() {

	suite.dbMock.On("DeleteProductById", mock.Anything, 1).Return(
		errors.New("Error while deleting Products"),
	)

	recorder := makeHTTPCall(http.MethodDelete,
		"/product/{product_id:[0-9]+}",
		"/product/1",
		"",
		deleteProductByIdHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), `{"error":{"message":"Internal server error  (Error feching data)"}}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *ProductsHandlerTestSuite) TestUpdateProductStockByIdSuccess() {

	suite.dbMock.On("GetProductByID", mock.Anything, 1).Return(testProduct, nil)
	suite.dbMock.On("UpdateProductStockById", mock.Anything, testProduct, 1).Return(testProduct, nil)

	recorder := makeHTTPCall(http.MethodPut,
		"/product/stock",
		"/product/stock?product_id=1&stock=0",
		"",
		updateProductStockByIdHandler(Dependencies{Store: suite.dbMock}),
	)
	//fmt.Println("recorder---->", recorder)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `{"data":{"id":1,"product_title":"test organization","description":"test@gmail.com","product_price":12,"discount":1,"tax":0.5,"stock":15,"category_id":5,"category":"2","brand":"IST","color":"Black","size":"Medium","image_urls":["url1","url2"]}}`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())

}

func (suite *ProductsHandlerTestSuite) TestUpdateProductStockByIdFailure() {

	suite.dbMock.On("UpdateProductStockById", mock.Anything, mock.Anything, "a").Return(db.Product{}, errors.New("Error id is missing/invalid"))

	recorder := makeHTTPCall(http.MethodPut,
		"/product/stock",
		"/product/stock?product_id=1&stock=a",
		"",
		updateProductStockByIdHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), `{"error":{"message":"Error id is missing/invalid"}}`, recorder.Body.String())

}
