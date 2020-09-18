package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	logger "github.com/sirupsen/logrus"
)

const (
	getProductCount = `SELECT count(id) from Products ;`
	getProductQuery = `SELECT id FROM products limit $1 OFFSET $2;`

	getProductIDQuery   = `SELECT id FROM products`
	getProductByIDQuery = `SELECT * FROM products WHERE id=$1`
	getCategoryByID     = `SELECT name FROM category WHERE id = $1`

	insertProductQuery = `INSERT INTO products ( name, description,
		  price, discount, tax, quantity, category_id, brand, color, size, image_url) VALUES ( 
		  :name, :description, :price, :discount, :tax, :quantity, :category_id, :brand, :color, :size, :image_url) RETURNING id;`
	deleteProductIdQuery    = `DELETE FROM products WHERE id = $1`
	updateProductStockQuery = `UPDATE products SET quantity= $1 where id = $2 `
	insertProductURLsQuery  = `INSERT INTO productimages (product_id, url) values ($1, $2)`
	updateProductQuery      = `UPDATE products SET name= $1, description=$2, price=$3, 
			discount=$4, tax=$5, quantity=$6, category_id=$7, brand=$8, color=$9, size=$10 WHERE id = $11`
)

type Product struct {
	Id           int            `db:"id" json:"id"`
	Name         string         `db:"name" json:"product_title"`
	Description  string         `db:"description" json:"description"`
	Price        float32        `db:"price" json:"product_price"`
	Discount     float32        `db:"discount" json:"discount"`
	Tax          float32        `db:"tax" json:"tax"`
	Quantity     int            `db:"quantity" json:"stock"`
	CategoryId   int            `db:"category_id" json:"category_id"`
	CategoryName string         `json:"category"`
	Brand        string         `db:"brand" json:"brand"`
	Color        string         `db:"color" json:"color"`
	Size         string         `db:"size" json:"size"`
	URLs         pq.StringArray `json:"image_url,omitempty" db:"image_url"`
}

// Pagination helps to return UI side with number of pages given a limitStr and pageStr number from Query Parameters
type Pagination struct {
	Products   []Product `json:"products"`
	TotalPages int       `json:"total_pages"`
}

func (product *Product) Validate() (map[string]ErrorResponse, bool) {
	var errorResponse map[string]ErrorResponse
	var valid bool

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
	if product.Discount < 0 || product.Discount > product.Price {
		fieldErrors["discount"] = "Can't be less than zero or more than Product's Price"
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
		return nil, valid
	}

	errorResponse = map[string]ErrorResponse{
		"error": ErrorResponse{
			Code:    "Invalid_data",
			Message: "Please Provide valid Product data",
			Fields:  fieldErrors,
		},
	}

	return errorResponse, false
}

// @Title GetProductByID
// @Description Get a Product Object by its Id
// @Params req.Context, product's Id
// @Returns Product Object, error if any
func (s *pgStore) GetProductByID(ctx context.Context, id int) (Product, error) {

	var product Product
	err := s.db.Get(&product, getProductByIDQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id: " + string(id))
		return Product{}, err
	}

	// Add category to Product's object
	var category string
	err = s.db.Get(&category, getCategoryByID, product.CategoryId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching category from database by product_id: " + string(id))
		return Product{}, err
	}

	product.CategoryName = category
	return product, nil
}

// @Title ListProducts
// @Description Get limited number of Products of particular pageStr
// @Params req.Context , limitStr, pageStr
// @Returns Count of Records, error if any
func (s *pgStore) ListProducts(ctx context.Context, limit int, page int) (int, []Product, error) {

	var count = 0
	var products []Product

	resultCount, err := s.db.Query(getProductCount)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Count of Products from database")
		return 0, []Product{}, err
	}

	if resultCount.Next() {
		err = resultCount.Scan(&count)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning Count Of Product into integer variable")
			return 0, []Product{}, err
		}
	}

	if count == 0 {
		err = fmt.Errorf("No records present")
		logger.WithField("err", err.Error()).Error("No Products were present in database")
		return 0, []Product{}, err
	}

	// error already handled in product_http
	//limit, _ := strconv.Atoi(limitStr)
	//page, _ := strconv.Atoi(pageStr)

	offset := (page - 1) * limit
	if (count - 1) < (limit * (page - 1)) {
		err = fmt.Errorf("Desired pageStr not found")
		logger.WithField("err", err.Error()).Error("pageStr Out Of range")
		return 0, []Product{}, err
	}

	result, err := s.db.Query(getProductQuery, limit, offset)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return 0, []Product{}, err
	}

	// idArr stores id's of all products
	var idArr []int

	for result.Next() {
		var id int
		err = result.Scan(&id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Couldn't Scan Product ids")
			return 0, []Product{}, err
		}
		idArr = append(idArr, id)
	}

	for i := 0; i < len(idArr); i++ {
		var product Product
		product, err = s.GetProductByID(ctx, int(idArr[i]))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error selecting Product from database by id " + string(idArr[i]))
			return 0, []Product{}, err
		}
		products = append(products, product)
	}

	return count, products, nil
}

func (s *pgStore) CreateProduct(ctx context.Context, product Product) (Product, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error beginning product insert transaction in db, CreateProduct with Name: " + product.Name)
		return Product{}, err
	}

	var row *sqlx.Rows
	row, err = tx.NamedQuery(insertProductQuery, product)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting product to database: " + product.Name)
		return Product{}, err
	}
	if row.Next() {
		err = row.Scan(&product.Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning product id from database: " + product.Name)
			return Product{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error commiting transaction inserting product into database: " + string(product.Name))
		return Product{}, err
	}

	return product, nil
}

func (s *pgStore) UpdateProductStockById(ctx context.Context, product Product, id int) (Product, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating update transaction")
		return Product{}, err
	}

	_, err = tx.Exec(updateProductStockQuery,
		product.Quantity,
		id,
	)
	if err != nil {
		// FAIL : Could not Update Product
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database Records not Found:" + string(id))
		return Product{}, err
	}

	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed. Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction updating product into database: " + string(id))
		return Product{}, err
	}

	return product, nil
}

func (s *pgStore) DeleteProductById(ctx context.Context, id int) error {

	rows, err := s.db.Exec(deleteProductIdQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting product" + string(id))
		return err
	}

	rows_affected, err := rows.RowsAffected()
	// if there is an error then roes_affected will by default be 0, so err != nil need not be handled separately
	if rows_affected == 0 {
		err = fmt.Errorf("Product doesn't exist in db, goodluck deleting it")
		return err
	}
	return nil
}

func (s *pgStore) UpdateProductById(ctx context.Context, product Product, id int) (Product, error) {

	var dbProduct Product
	err := s.db.Get(&dbProduct, getProductByIDQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching product ")
		return Product{}, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating update Product transaction")
		return Product{}, err
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
	if product.CategoryId == 0 {
		product.CategoryId = dbProduct.CategoryId
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

	_, valid := product.Validate()
	if !valid {
		return Product{}, fmt.Errorf("Product Validation failed. Invalid Fields present in the product e.g Discount is greater than Price")
	}

	_, err = tx.Exec(updateProductQuery,
		product.Name,
		product.Description,
		product.Price,
		product.Discount,
		product.Tax,
		product.Quantity,
		product.CategoryId,
		product.Brand,
		product.Color,
		product.Size,
		id,
	)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database :" + string(id))
		return Product{}, err
	}

	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed. Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction updating product into database: " + string(id))
		return Product{}, nil
	}

	return product, nil
}
