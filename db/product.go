package db

import (
	"context"
	"fmt"
	"strconv"

	//"database/sql"
	logger "github.com/sirupsen/logrus"
)

const (
	getProductCount = `SELECT count(id) from Products ;`
	getProductQuery = `SELECT id FROM products LIMIT $1  OFFSET  ($2 -1) * $1;`

	getProductIDQuery     = `SELECT id FROM products`
	getProductByIDQuery   = `SELECT * FROM products WHERE id=$1`
	getProductByNameQuery = `SELECT * FROM products WHERE name=$1`
	getCategoryByID       = `SELECT name FROM category WHERE id = $1`
	insertProductQuery    = `INSERT INTO products ( name, description,
		  price, discount, tax, quantity, category_id, brand, color, size) VALUES ( 
		  :name, :description, :price, :discount, :tax, :quantity, :category_id, :brand, :color, :size)`
	deleteProductIdQuery    = `DELETE FROM products WHERE id = $1`
	updateProductStockQuery = `UPDATE products SET quantity= $1 where id = $2`
	newInsertRecord         = `SELECT MAX(id) from products`
	insertProductURLsQuery  = `INSERT INTO productimages (product_id, url) values ($1, $2)`
	updateProductQuery      = `UPDATE products SET name= $1, description=$2, price=$3, 
			discount=$4, tax=$5, quantity=$6, category_id=$7, brand=$8, color=$9, size=$10 WHERE id = $11`
	updateProductImageQuery = `UPDATE productimages SET url = $1 WHERE product_id = $2`
)

type Product struct {
	ID           int      `db:"id" json:"id"`
	Name         string   `db:"name" json:"product_title"`
	Description  string   `db:"description" json:"description"`
	Price        float32  `db:"price" json:"product_price"`
	Discount     float32  `db:"discount" json:"discount"`
	Tax          float32  `db:"tax" json:"tax"`
	Quantity     int      `db:"quantity" json:"stock"`
	CategoryID   int      `db:"category_id" json:"category_id"`
	CategoryName string   `json:"category"`
	Brand        string   `db:"brand" json:"brand"`
	Color        string   `db:"color" json:"color"`
	Size         string   `db:"size" json:"size"`
	URLs         []string `json:"image_url,omitempty"`
}

// Pagination helps to return UI side with number of pages given a limit and page
type Pagination struct {
	Products   []Product `json:"products"`
	TotalPages int       `json:"total_pages"`
}

func (product *Product) Validate() (errorResponse map[string]ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if product.Name == "" {
		fieldErrors["product_name"] = "Can't be blank"
	}
	if product.Description == "" {
		fieldErrors["product_description"] = "Can't be blank "
	}
	if product.Price <= 0 {
		fieldErrors["price"] = "Can't be blank  or less than zero"
	}
	if product.Discount < 0 {
		fieldErrors["discount"] = "Can't be less than zero"
	}
	if product.Tax < 0 {
		fieldErrors["tax"] = "Can't be less than zero"
	}
	// If Quantity gets's < 0 by UpdateProductStockById Method, this is what saves us
	if product.Quantity < 0 {
		fieldErrors["available_quantity"] = "Can't be blank or less than zero"
	}
	if product.CategoryID == 0 {
		fieldErrors["category_id"] = "Can't be blank"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]ErrorResponse{
		"error": ErrorResponse{
			Code:    "Invalid_data",
			Message: "Please Provide valid Product data",
			Fields:  fieldErrors,
		},
	}
	// TODO Other Validations
	return
}

func (product *Product) PartialValidate() (errorResponse map[string]ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if product.Price < 0 {
		fieldErrors["price"] = "Can't be blank  or less than zero"
	}
	if product.Discount < 0 {
		fieldErrors["discount"] = "Can't be less than zero"
	}
	if product.Tax < 0 {
		fieldErrors["tax"] = "Can't be less than zero"
	}
	// If Quantity gets's < 0 by UpdateProductStockById Method, this is what saves us
	if product.Quantity < 0 {
		fieldErrors["available_quantity"] = "Can't be blank or less than zero"
	}
	if product.CategoryID < 0 {
		fieldErrors["category_id"] = "Can't be invalid"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]ErrorResponse{
		"error": ErrorResponse{
			Code:    "Invalid_data",
			Message: "Please Provide valid Product data",
			Fields:  fieldErrors,
		},
	}
	// TODO Other Validations
	return
}

// @Title GetProductByID
// @Description Get a Product Object by its Id
// @Params req.Context, product's Id
// @Returns Product Object, error if any
func (s *pgStore) GetProductByID(ctx context.Context, Id int) (product Product, err error) {

	err = s.db.Get(&product, getProductByIDQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id: " + string(Id))
		return
	}

	// Add category to Product's object
	var category string
	err = s.db.Get(&category, getCategoryByID, product.CategoryID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching category from database by product_id: " + string(Id))
		return
	}

	product.CategoryName = category

	// Add Product Image URL's to Product's Object
	productImage, err := s.GetProductImagesByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting productImage from database by id " + string(Id))
		return
	}

	for j := 0; j < len(productImage); j++ {
		product.URLs = append(product.URLs, productImage[j].URL)
	}

	return
}

// @Title ListProducts
// @Description Get limited number of Products of particular page
// @Params req.Context , limit, page
// @Returns Count of Records, error if any
func (s *pgStore) ListProducts(ctx context.Context, limit string, page string) (count int, products []Product, err error) {

	resultCount, err := s.db.Query(getProductCount)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Count of Products from database")
		return
	}

	for resultCount.Next() {
		err = resultCount.Scan(&count)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning Count Of Product into integer variable")
			return
		}
	}

	if count == 0 {
		err = fmt.Errorf("No records present")
		logger.WithField("err", err.Error()).Error("No Products were present in database")
		return
	}

	// error already handled in product_http
	ls, _ := strconv.Atoi(limit)
	ps, _ := strconv.Atoi(page)

	if (count - 1) < (int(ls) * (int(ps) - 1)) {
		err = fmt.Errorf("Desired Page not found")
		logger.WithField("err", err.Error()).Error("Page Out Of range")
		return
	}

	result, err := s.db.Query(getProductQuery, ls, ps)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

	// idArr stores id's of all products
	var idArr []int

	for result.Next() {
		var Id int
		err = result.Scan(&Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Couldn't Scan Product ids")
			return
		}
		idArr = append(idArr, Id)
	}

	for i := 0; i < len(idArr); i++ {
		var product Product
		product, err = s.GetProductByID(ctx, int(idArr[i]))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error selecting Product from database by id " + string(idArr[i]))
			return
		}
		products = append(products, product)
	}

	return
}

// CreateNewProduct
func (s *pgStore) CreateNewProduct(ctx context.Context, p Product) (createdProduct Product, err error) {
	// First, make sure Product isn't already in db, if Product is present, just return the it
	err = s.db.Get(&createdProduct, getProductByNameQuery, p.Name)
	if err == nil {
		// If there's already a product, err wil be nil, so no new Product is populated.
		err = fmt.Errorf("Product Already exists!")
		return
	}
	tx, err := s.db.Beginx() // Use Beginx instead of MustBegin so process doesn't die if there is an error
	if err != nil {
		// FAIL : Could not begin database transaction
		logger.WithField("err", err.Error()).Error("Error beginning product insert transaction in db.CreateNewProduct with Id: " + string(p.ID))
		return
	}

	_, err = tx.NamedExec(insertProductQuery, p)
	//  p.Name, p.Description, p.Price, p.Discount, p.Tax, p.Quantity, p.CategoryId, p.Brand, p.Color, p.Size

	if err != nil {
		// FAIL : Could not run insert Query
		logger.WithField("err", err.Error()).Error("Error inserting product to database: " + p.Name)
		return
	}
	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed.Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction inserting product into database: " + string(p.ID))
		return
	}

	//length of url
	urls := len(p.URLs)
	var number int
	//new insert record get id number
	result, err := s.db.Query(newInsertRecord)
	for result.Next() {
		err = result.Scan(&number)
	}

	for i := 0; i < urls; i++ {
		//insert urls of given records
		_, err = s.db.Exec(insertProductURLsQuery, number, p.URLs[i])
		if err != nil {
			// FAIL : Could not run insert Query
			logger.WithField("err", err.Error()).Error("Error inserting urls to database: ")
			return
		}
	}

	// Re-select Product and return it
	createdProduct, err = s.GetProductByID(ctx, number)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting from database with id: " + string(p.ID))
		return
	}
	return
}

func (s *pgStore) UpdateProductStockById(ctx context.Context, product Product, Id int) (updatedProduct Product, err error) {

	var dbProduct Product
	err = s.db.Get(&dbProduct, getProductByIDQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching product ")
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating update transaction")
		return
	}

	_, err = tx.Exec(updateProductStockQuery,
		product.Quantity,
		Id,
	)

	if err != nil {
		// FAIL : Could not Update Product
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database :" + string(Id))
		return
	}

	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed. Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction updating product into database: " + string(Id))
		return
	}

	updatedProduct, err = s.GetProductByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while getting updated product ")
		return
	}
	return
}

func (s *pgStore) DeleteProductById(ctx context.Context, Id int) (err error) {

	rows, err := s.db.Exec(deleteProductIdQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting product" + string(Id))
		return
	}

	rows_affected, err := rows.RowsAffected()
	if rows_affected == 0 {
		err = fmt.Errorf("Product doesn't exist in db, goodluck deleting it")
	}
	return
}

func (s *pgStore) UpdateProductById(ctx context.Context, product Product, Id int) (updatedProduct Product, err error) {

	var dbProduct Product
	err = s.db.Get(&dbProduct, getProductByIDQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching product ")
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating update Product transaction")
		return
	}

	if product.Name == "" {
		product.Name = dbProduct.Name
	}
	if product.Description == "" {
		product.Description = dbProduct.Description
	}
	if product.Price == 0 {
		product.Price = dbProduct.Price
	}
	if product.Discount == 0 {
		product.Discount = dbProduct.Discount
	}
	if product.Tax == 0 {
		product.Tax = dbProduct.Tax
	}
	if product.Quantity == 0 {
		product.Quantity = dbProduct.Quantity
	}
	if product.CategoryID == 0 {
		product.CategoryID = dbProduct.CategoryID
	}
	if product.Brand == "" {
		product.Brand = dbProduct.Brand
	}
	if product.Color == "" {
		product.Color = dbProduct.Color
	}
	if product.Size == "" {
		product.Size = dbProduct.Size
	}

	_, err = tx.Exec(updateProductQuery,
		product.Name,
		product.Description,
		product.Price,
		product.Discount,
		product.Tax,
		product.Quantity,
		product.CategoryID,
		product.Brand,
		product.Color,
		product.Size,
		Id,
	)

	if len(product.URLs) != 0 {
		_, err = tx.Exec(updateProductImageQuery, product.URLs[0], Id)

	}

	if err != nil {
		// FAIL : Could not Update Product
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database :" + string(Id))
		return
	}

	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed. Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction updating product into database: " + string(Id))
		return
	}

	updatedProduct, err = s.GetProductByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while getting updated product ")
		return
	}
	return
}
