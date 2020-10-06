package db

import (
	"context"
	"github.com/lib/pq"
	logger "github.com/sirupsen/logrus"
)

const (
	joinCartProductQuery = `SELECT cart.product_id, products.name as product_name, cart.quantity, category.cname as category_name, products.price , products.description, products.image_urls, products.discount, products.tax from cart inner join products on cart.product_id=products.id inner join category on category.cid=products.cid where cart.id=$1 ORDER BY cart.product_id ASC;`
)

// CartProduct is rquirement of frontend as json response
type CartProduct struct {
	Id          int            `db:"product_id" json:"id"`
	Name        string         `db:"product_name" json:"product_title"`
	Quantity    int            `db:"quantity" json:"quantity"`
	Category    string         `db:"category_name" json:"category,omitempty"`
	Price       float64        `db:"price" json:"product_price"`
	Description string         `db:"description" json:"description"`
	ImageUrls   pq.StringArray `db:"image_urls" json:"image_url,*"`
	Discount    float32        `db:"discount" json:"discount"`
	Tax         float32        `db:"tax" json:"tax"`
}

// *pgStore  because deps.Store.GetCart - deps is of struct Dependencies - vch is assigned to db conn obj
func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	err = s.db.Select(&cart_products, joinCartProductQuery, user_id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching data from cart")
		return
	}
	return
}

func (s *pgStore) AddToCart(ctx context.Context, cartID, productID int) (rowsAffected int64, err error) {
	insert := `INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)`
	result, err := s.db.Exec(insert, cartID, productID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error adding to cart")
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching affected rows")
	}
	return
}

func (s *pgStore) DeleteFromCart(ctx context.Context, cartID, productID int) (rowsAffected int64, err error) {
	delete := `DELETE FROM cart WHERE id = $1 AND product_id = $2`
	result, err := s.db.Exec(delete, cartID, productID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while removing from cart")
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching affected rows")
	}
	return
}

func (s *pgStore) UpdateIntoCart(ctx context.Context, quantity, cartID, productID int) (rowsAffected int64, err error) {
	update := `UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3`
	result, err := s.db.Exec(update, quantity, cartID, productID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while updating into cart")
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching affected rows")
	}
	return
}
