package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

// CartProduct is rquirement of frontend as json response
type CartProduct struct {
	Id          int     `db:"product_id" json:"id"`
	Name        string  `db:"product_name" json:"product_title"`
	Quantity    int     `db:"quantity" json:"quantity"`
	Category    string  `db:"category_name" json:"category,omitempty"`
	Price       float64 `db:"price" json:"product_price"`
	Description string  `db:"description" json:"description"`
	ImageUrls   string  `db:"url" json:"image_url,omitempty"`
}

const (
	joinCartProductQuery = `SELECT cart.product_id, products.name as product_name, cart.quantity, category.name as category_name, products.price , products.description, productimages.url from cart inner join products on cart.product_id=products.id inner join category on category.id=products.category_id inner join productimages on products.id=productimages.product_id where cart.id=$1 ORDER BY cart.product_id ASC;`
)

// *pgStore  because deps.Store.GetCart - deps is of struct Dependencies - vch is assigned to db conn obj
func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	err = s.db.Select(&cart_products, joinCartProductQuery, user_id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching data from cart")
		return
	}
	return
}
