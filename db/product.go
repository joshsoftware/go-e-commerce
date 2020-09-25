package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	logger "github.com/sirupsen/logrus"
)

const (
	getProductCount     = `SELECT count(id) from Products ;`
	getProductQuery     = `SELECT * FROM products p INNER JOIN category c ON p.cid = c.cid ORDER BY p.id LIMIT $1 OFFSET $2 ;`
	getProductByIDQuery = `SELECT * FROM products p INNER JOIN category c ON p.cid = c.cid WHERE p.id=$1`
	insertProductQuery  = `INSERT INTO products ( name, description,
		  price, discount, tax, quantity, cid, brand, color, size, image_urls) VALUES ( 
		  :name, :description, :price, :discount, :tax, :quantity, :cid, :brand, :color, :size, :image_urls) RETURNING id;`
	deleteProductIdQuery    = `DELETE FROM products WHERE id = $1`
	updateProductStockQuery = `UPDATE products SET quantity= $1 where id = $2 `
	updateProductQuery      = `UPDATE products SET name= $1, description=$2, price=$3, 
			discount=$4, tax=$5, quantity=$6, cid=$7, brand=$8, color=$9, size=$10, image_urls=$11 WHERE id = $12`
)

type Product struct {
	Id           int            `db:"id" json:"id" schema:"-"`
	Name         string         `db:"name" json:"product_title" schema:"product_title"`
	Description  string         `db:"description" json:"description" schema:"description"`
	Price        float32        `db:"price" json:"product_price" schema:"product_price"`
	Discount     float32        `db:"discount" json:"discount" schema:"discount"`
	Tax          float32        `db:"tax" json:"tax" schema:"tax"`
	Quantity     int            `db:"quantity" json:"stock" schema:"stock"`
	CategoryId   int            `db:"cid" json:"category_id" schema:"category_id"`
	CategoryName string         `db:"cname" json:"category" schema:"category"`
	Brand        string         `db:"brand" json:"brand" schema:"brand"`
	Color        *string        `db:"color" json:"color,*" schema:"color,*"`
	Size         *string        `db:"size" json:"size,*" schema:"size,*"`
	URLs         pq.StringArray `db:"image_urls" json:"image_urls,*"  schema:"-"`
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
	return product, nil
}

// @Title ListProducts
// @Description Get limited number of Products of particular pageStr
// @Params req.Context , limitStr, pageStr
// @Returns Count of Records, error if any
func (s *pgStore) ListProducts(ctx context.Context, limit int, offset int) (int, []Product, error) {

	var totalRecords int
	var products []Product

	resultCount, err := s.db.Query(getProductCount)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Count of Products from database")
		return 0, []Product{}, err
	}

	if resultCount.Next() {
		err = resultCount.Scan(&totalRecords)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning Count Of Product into integer variable")
			return 0, []Product{}, err
		}
	}

	if totalRecords-1 < offset {
		err = fmt.Errorf("Page out of Range!")
		logger.WithField("err", err.Error()).Error("Error Offset is greater than total records")
		return 0, []Product{}, err

	}

	err = s.db.Select(&products, getProductQuery, limit, offset)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return 0, []Product{}, err
	}

	if products == nil {
		err = fmt.Errorf("Desired page not found")
		logger.WithField("err", err.Error()).Error("page Out Of range")
		return 0, []Product{}, err
	}

	return totalRecords, products, nil
}

func (s *pgStore) CreateProduct(ctx context.Context, product Product) (int, error) {

	var row *sqlx.Rows
	row, err := s.db.NamedQuery(insertProductQuery, product)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting product to database: " + product.Name)
		return 0, err
	}
	if row.Next() {
		err = row.Scan(&product.Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning product id from database: " + product.Name)
			return 0, err
		}
	}
	row.Close()

	return product.Id, nil
}

func (s *pgStore) UpdateProductStockById(ctx context.Context, product Product, id int) (Product, error) {

	_, err := s.db.Exec(updateProductStockQuery,
		product.Quantity,
		id,
	)
	if err != nil {
		// FAIL : Could not Update Product
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database Records not Found:" + string(id))
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

	if product.URLs != nil {
		//fmt.Println("Db product name--->", dbProduct.Name)
		var files []string

		root := "./assets/productImages/"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			logger.WithField("err", err.Error()).Error("cannot fetch products images path from assets/productImages dir for updation")
			return Product{}, err
		}
		for _, file := range files {
			if strings.Index(file, dbProduct.Name) > 0 {
				err := os.Remove(file)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
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
	if product.Color == nil || *product.Color == "" {
		*product.Color = *dbProduct.Color
	}
	if product.Size == nil || *product.Color == "" {
		*product.Size = *dbProduct.Size
	}
	if product.URLs == nil {
		product.URLs = dbProduct.URLs
	}

	_, valid := product.Validate()
	if !valid {
		return Product{}, fmt.Errorf("Product Validation failed. Invalid Fields present in the product e.g Discount is greater than Price")
	}

	_, err = s.db.Exec(updateProductQuery,
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
		product.URLs,
		id,
	)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database :" + string(id))
		return Product{}, err
	}

	product.Id = id

	return product, nil
}
