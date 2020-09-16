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

type Products struct {
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Price       float32 `db:"price" json:"price"`
	Discount    float32 `db:"discount" json:"discount"`
	Tax         float32 `db:"tax" json:""tax`
	Quantity    int     `db:"quantity" json:"quantity"`
	CategoryId  int     `db:"category_id" json:"category_id"`
}

// CartProduct is rquirement of frontend as json response
type CartProduct struct {
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"product_title" json:"product_title"`
	Description string  `db:"description" json:"description"`
	Quantity    int     `db:"quantity" json:"quantity"`
	Price       float32 `db:"price" json:"product_price"`
	Discount    float32 `db:"discount" json:"discount"`
	Tax         float32 `db:"tax" json:"tax"`
	Category    string  `db:"category" json:"category"`
	ImageUrl    string  `db:"image_url" json:"image_url"`
}

const (
	getCartQuery         = `SELECT product_id  FROM cart WHERE id=$1`
	getCartQuantityQuery = `SELECT quantity FROM cart WHERE id=$1`
	getProductsQuery     = `SELECT id,name,description,price,discount,tax,quantity,category_id FROM products WHERE id IN (?)`
	getCategoryQuery     = `SELECT name from category where id=$1`
	getProductImageQuery = `SELECT url from productimages where product_id=$1`
)

func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart_products []CartProduct, err error) {
	var pids []interface{}
	var quantities []int
	var category, image_url []string
	var products []Products

	err = s.db.Select(&pids, getCartQuery, user_id)
	err = s.db.Select(&quantities, getCartQuantityQuery, user_id)
	// fmt.printf("pids : %v  quantities :%v", pids, quantities)
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

	for _, product := range products {
		err = s.db.Select(&category, getCategoryQuery, product.CategoryId)
		err = s.db.Select(&image_url, getProductImageQuery, product.Id)
	}

	for index, product := range products {
		cart_products = append(
			cart_products,
			CartProduct{
				Id:          product.Id,
				Quantity:    quantities[index],
				Category:    category[index],
				Price:       product.Price,
				Description: product.Description,
				ImageUrl:    image_url[index],
				Name:        product.Name,
				Discount:    product.Discount,
				Tax:         product.Tax,
			})
	}

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing cart")
		return
	}
	return
}

func (s *pgStore) AddToCart(ctx context.Context, cartId, productId int) (rowsAffected int64, err error) {
	insert := `INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)`
	result, err := s.db.Exec(insert, cartId, productId)

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

func (s *pgStore) RemoveFromCart(ctx context.Context, cartId, productId int) (rowsAffected int64, err error) {
	delete := `DELETE FROM cart WHERE id = $1 AND product_id = $2`
	result, err := s.db.Exec(delete, cartId, productId)

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

func (s *pgStore) UpdateIntoCart(ctx context.Context, quantity, cartId, productId int) (rowsAffected int64, err error) {
	update := `UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3`
	result, err := s.db.Exec(update, quantity, cartId, productId)

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
