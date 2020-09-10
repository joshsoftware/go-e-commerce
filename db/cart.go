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
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

const (
	getCartQuery     = `SELECT product_id FROM cart WHERE id=$1`
	getProductsQuery = `SELECT * FROM products WHERE id IN (?)`
)

func (s *pgStore) GetCart(ctx context.Context, user_id int) (products []Product, err error) {
	var pids []interface{}
	err = s.db.Select(&pids, getCartQuery, user_id)
	query, args, err := sqlx.In(getProductsQuery, pids)
	query = s.db.Rebind(query)
	err = s.db.Select(&products, query, args...)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing cart")
		return
	}
	return
}
