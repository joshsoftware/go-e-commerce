package db

import (
	"context"
	"fmt"

	//"database/sql"
	logger "github.com/sirupsen/logrus"
)

const (
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
)

type Product struct {
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"product_title"`
	Description string  `db:"description" json:"description"`
	Price       float32 `db:"price" json:"product_price"`
	Discount    float32 `db:"discount" json:"discount"`
	Tax         float32 `db:"tax" json:"tax"`
	Quantity    int     `db:"quantity" json:"stock"`
	CategoryId  int     `db:"category_id" json:"category_id"`
	// CategoryName doesn't exists in Product db,
	// it is fetched from category table whenever
	// products are to be returned
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
	if product.CategoryId == 0 {
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
	err = s.db.Get(&category, getCategoryByID, product.CategoryId)
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

// @Title TotalRecords
// @Description Get count of Total Product Records
// @Params req.Context
// @Returns Count of Records, error if any
func (s *pgStore) TotalRecords(ctx context.Context) (total int, err error) {

	getTotalRecord := "SELECT count(id) from products ;"
	result, err := s.db.Query(getTotalRecord)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

	for result.Next() {
		err = result.Scan(&total)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning Product Ids into integer variable")
			return
		}
	}
	return
}

func (s *pgStore) TotalSearchproduct(ctx context.Context, search string) (total int, err error) {

	getTotalRecord := `SELECT count(p.id) from products p
			 INNER JOIN category c 
			 ON p.category_id = c.id
			 WHERE`

	getTotalRecord += " p.name LIKE '%" + search + "%' OR p.description  LIKE '%" + search + "%' OR p.brand  LIKE  '%" + search + "%' OR p.color LIKE  '%" + search + "%' OR p.size LIKE  '%" + search + "%' OR c.name LIKE  '%" + search + "%' OR c.description LIKE '%" + search + "%';"

	fmt.Println("GIven Query--->", getTotalRecord)
	result, err := s.db.Query(getTotalRecord)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

	for result.Next() {
		err = result.Scan(&total)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning Product Ids into integer variable")
			return
		}
	}
	return
}

func (s *pgStore) SerachProducts(ctx context.Context, search string, limit string, page string) (products []Product, err error) {

	var idArr []int
	getTotalRecord := `SELECT p.id from products p
			 INNER JOIN category c 
			 ON p.category_id = c.id
			 WHERE`

	getTotalRecord += " p.name LIKE '%" + search + "%' OR p.description  LIKE '%" + search + "%' OR p.brand  LIKE  '%" + search + "%' OR p.color LIKE  '%" + search + "%' OR p.size LIKE  '%" + search + "%' OR c.name LIKE  '%" + search + "%' OR c.description LIKE '%" + search + "%'"
	getTotalRecord += " LIMIT " + string(limit) + "  OFFSET  (" + string(page) + " -1) * " + string(limit) + ""
	getTotalRecord += " ;"

	result, err := s.db.Query(getTotalRecord)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

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

// @Title ListProducts
// @Description Get limited number of Products of particular page
// @Params req.Context , limit, page
// @Returns Count of Records, error if any
func (s *pgStore) ListProducts(ctx context.Context, limit string, page string) (products []Product, err error) {

	// idArr stores id's of all products
	var idArr []int

	getProductQuery := `SELECT id FROM products`
	getProductQuery += " LIMIT " + string(limit) + "  OFFSET  (" + string(page) + " -1) * " + string(limit) + ""
	getProductQuery += " ;"

	result, err := s.db.Query(getProductQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

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
		logger.WithField("err", err.Error()).Error("Error beginning product insert transaction in db.CreateNewProduct with Id: " + string(p.Id))
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
		logger.WithField("err", err.Error()).Error("Error commiting transaction inserting product into database: " + string(p.Id))
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
		logger.WithField("err", err.Error()).Error("Error selecting from database with id: " + string(p.Id))
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
