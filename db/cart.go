package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

type Cart struct {
	Id        int `db:"id" json:"id"`
	ProductId int `db:"product_id" json:"product_id"`
	Quantity  int `db:"quantity" json:"quantity"`
}

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
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"product_title" json:"product_title"`
	Quantity    int     `db:"quantity" json:"quantity"`
	Category    string  `db:"category" json:"category"`
	Price       float64 `db:"price" json:"product_price"`
	Description string  `db:"description" json:"description"`
	ImageUrl    string  `db:"image_url" json:"image_url"`
}

const (
	getCartQuery         = `SELECT product_id FROM cart WHERE id=$1`
	getProductsQuery     = `SELECT * FROM products WHERE id IN (?)`
	getCategoryQuery     = `SELECT name from category where id=$1`
	getProductImageQuery = `SELECT url from productimages where id=$1`
)

func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	var pids []interface{}
	var category, image_url []string
	var products []Product

	err = s.db.Select(&pids, getCartQuery, user_id)

	query, args, err := sqlx.In(getProductsQuery, pids)
	query = s.db.Rebind(query)
	err = s.db.Select(&products, query, args...)

	for _, product := range products {
		err = s.db.Select(&category, getCategoryQuery, product.CategoryId)
		err = s.db.Select(&image_url, getProductImageQuery, product.Id)
	}

	for index, product := range products {
		cart_products = append(
			cart_products,
			CartProduct{
				Id:       product.Id,
				Quantity: product.Quantity,
				Category: category[index],
				Price:    product.Price, Description: product.Description,
				ImageUrl: image_url[index],
				Name:     product.Name})
	}

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing cart")
		return
	}
	return
}
