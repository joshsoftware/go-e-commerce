package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

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

type JoinCartProduct struct {
	ProductId   int     `db:"product_id"`
	Name        string  `db:"name"`
	Quantity    int     `db:"quantity"`
	CategoryId  int     `db:"category_id"`
	Price       float64 `db:"price"`
	Description string  `db:"description"`
}

const (
	joinCartProductQuery = `SELECT cart.product_id, products.name, cart.quantity, products.category_id, products.price , products.description from cart inner join products on cart.product_id=products.id where cart.id=$1;`
	getCategoryQuery     = `SELECT name from category where id=$1`
	getProductImageQuery = `SELECT url from productimages where product_id=$1`
)

func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	var joinCartProduct []JoinCartProduct

	err = s.db.Select(&joinCartProduct, joinCartProductQuery, user_id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching data from cart")
		return
	}

	for _, row := range joinCartProduct {
		var (
			category_name []string
			image_urls    []string
		)
		err = s.db.Select(&category_name, getCategoryQuery, row.CategoryId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data from cart")
			return
		}

		err = s.db.Select(&image_urls, getProductImageQuery, row.ProductId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data from cart")
			return
		}

		cart_products = append(
			cart_products,
			CartProduct{
				Id:          row.ProductId,
				Quantity:    row.Quantity,
				Category:    category_name[0],
				Price:       row.Price,
				Description: row.Description,
				ImageUrls:   image_urls,
				Name:        row.Name,
			})

	}
	return
}
