package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

type Cart struct {
	Id        int    `db:"id" json:"id"`
	ProductId string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
}

const (
	getCartQuery = `SELECT * FROM cart WHERE id=$1`
)

func (s *pgStore) GetCart(ctx context.Context, user_id int) (cart []Cart, err error) {
	err = s.db.Select(&cart, getCartQuery, user_id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing cart")
		return
	}

	return
}
