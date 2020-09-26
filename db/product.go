package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"regexp"

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
	deleteProductIdQuery    = `DELETE FROM products WHERE id = $1 RETURNING image_urls`
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
	Color        string         `db:"color" json:"color,*" schema:"color,*"`
	Size         string         `db:"size" json:"size,*" schema:"size,*"`
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

func deleteImages(files pq.StringArray) error {

	root := "./assets/productImages/"
	for _, file := range files {
		file = root + file
		fmt.Println(file)
		err := os.Remove(file)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error Couldn't remove the file!")
			return err
		}
	}
	return nil
}

func imagesStore(images []*multipart.FileHeader, product *Product) error {

	for i, _ := range images {
		image, err := images[i].Open()
		defer image.Close()
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding image Data, probably invalid image")
			return err
		}

		extensionRegex := regexp.MustCompile(`[.]+.*`)
		extension := extensionRegex.Find([]byte(images[i].Filename))
		// normally our extensions be like .jpg, .jpeg, .png etc
		if len(extension) < 2 || len(extension) > 5 {
			err = fmt.Errorf("Couldn't get extension of file!")
			logger.WithField("err", err.Error()).Error("Error while getting image Extension. Re-check the image file extension!")

			return err
		}

		directoryPath := "assets/productImages"

		tempFile, err := ioutil.TempFile(directoryPath, (*product).Name+"-*"+string(extension))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while Creating a Temporary File")
			return err
		}
		defer tempFile.Close()

		imageBytes, err := ioutil.ReadAll(image)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while reading image File")
			return err
		}
		tempFile.Write(imageBytes)
		(*product).URLs = append(product.URLs, tempFile.Name()[len(directoryPath)+1:])
	}
	return nil
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

func (s *pgStore) CreateProduct(ctx context.Context, product Product, images []*multipart.FileHeader) (Product, error) {

	if images != nil {
		err := imagesStore(images, &product)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error inserting images in assets: " + product.Name)
			return Product{}, err
		}
	}

	var row *sqlx.Rows
	row, err := s.db.NamedQuery(insertProductQuery, product)
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
	row.Close()
	return product, nil
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

	var files pq.StringArray
	result, err := s.db.Queryx(deleteProductIdQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting product" + string(id))
		return err
	}

	if result.Next() {
		err = result.Scan(&files)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning image_urls into files variable")
			return err
		}
	} else {
		err = fmt.Errorf("Product doesn't exist in db, goodluck deleting it")
		return err
	}

	// do not throw error as deletion of Product data was successful.
	if files != nil {
		err = deleteImages(files)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error Couldn't remove the images file!")
		}
	}

	return nil
}

func (s *pgStore) UpdateProductById(ctx context.Context, product Product, id int, images []*multipart.FileHeader) (Product, error) {

	var dbProduct Product
	err := s.db.Get(&dbProduct, getProductByIDQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching product ")
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
	if product.CategoryName == "" {
		product.CategoryName = dbProduct.CategoryName
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

	if images != nil {
		err = imagesStore(images, &product)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error inserting images in assets: " + product.Name)
			return Product{}, err
		}
	}

	// Update images only after validations
	if product.URLs != nil {
		files := dbProduct.URLs
		err = deleteImages(files)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error Couldn't remove the images file!")
			return Product{}, err
		}
	} else {
		product.URLs = dbProduct.URLs
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
