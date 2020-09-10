package db

import (
	"context"

	//"database/sql"
	logger "github.com/sirupsen/logrus"
)

const (
	getProductIDQuery = `SELECT id FROM products`
	// id is PRIMARY KEY, so no need to limit
	getProductByIDQuery = `SELECT * FROM products WHERE id=$1`

	getCategoryByID = `SELECT name FROM category WHERE id = $1`
)

type Product struct {
	Id           int      `db:"id" json:"product_id"`
	Name         string   `db:"name" json:"product_name"`
	Description  string   `db:"description" json:"product_description"`
	Price        float32  `db:"price" json:"price"`
	Discount     float32  `db:"discount" json:"discount"`
	Quantity     int      `db:"quantity" json:"available_quantity"`
	CategoryId   int      `db:"category_id" json:"category_id"`
	CategoryName string   `json:"category_name,omitempty"`
	URLs         []string `json:"productimage_urls,omitempty"`
}

func (s *pgStore) GetProductByID(ctx context.Context, Id int) (product Product, err error) {

	err = s.db.Get(&product, getProductByIDQuery, Id)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id: " + string(Id))
		return
	}

	var category string
	err = s.db.Get(&category, getCategoryByID, product.CategoryId)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching category from database by product_id: " + string(Id))
		return
	}

	product.CategoryName = category

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

func (s *pgStore) ListProducts(ctx context.Context) (products []Product, err error) {

	// idArr stores id's of all products
	var idArr []int

	result, err := s.db.Query(getProductIDQuery)
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
