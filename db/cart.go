package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

type Product struct {
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
	Discount    float64 `db:"discount" json:"discount"`
	Quantity    int     `db:"quantity" json:"quantity"`
	CategoryId  int     `db:"category_id" json:"category_id"`
}

// CartProduct is rquirement of frontend as json response
type CartProduct struct {
	Id          int      `db:"id" json:"id"`
	Name        string   `db:"product_title" json:"product_title"`
	Quantity    int      `db:"quantity" json:"quantity"`
	Category    string   `db:"category" json:"category,omitempty"`
	Price       float64  `db:"price" json:"product_price"`
	Description string   `db:"description" json:"description"`
	ImageUrls   []string `db:"image_url" json:"image_url,omitempty"`
}

const (
	getCartQuery         = `SELECT product_id FROM cart WHERE id=$1`
	geCarttQuantityQuery = `SELECT quantity FROM cart WHERE id=$1`
	getProductsQuery     = `SELECT * FROM products WHERE id IN (?)`
	getCategoryQuery     = `SELECT name from category where id=$1`
	getProductImageQuery = `SELECT url from productimages where product_id=$1`
)

func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	var pids, quantities []int
	var products []Product

	err = s.db.Select(&pids, getCartQuery, user_id)
	err = s.db.Select(&quantities, geCarttQuantityQuery, user_id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching data from cart")
		return
	}

	query, args, err := sqlx.In(getProductsQuery, pids)
	query = s.db.Rebind(query)

	err = s.db.Select(&products, query, args...)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching data from cart")
		return
	}

	for index, product := range products {
		var category, image_urls []string
		err = s.db.Select(&category, getCategoryQuery, product.CategoryId)
		err = s.db.Select(&image_urls, getProductImageQuery, product.Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error listing cart")
			return
		}
		fmt.Println(category, image_urls, index)
		cart_products = append(
			cart_products,
			CartProduct{
				Id:          product.Id,
				Quantity:    quantities[index],
				Category:    category[0],
				Price:       product.Price,
				Description: product.Description,
				ImageUrls:   image_urls,
				Name:        product.Name})

	}
	return
}
