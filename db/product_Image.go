package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

const (
	getProductImagesByIDQuery = `SELECT * FROM productimages WHERE product_id = $1 ORDER BY url ASC`
)

type ProductImage struct {
	ProductId   int    `db:"product_id" json:"product_id"`
	URL         string `db:"url" json:"url"`
	Description string `db:"description" json:"image_description"`
}

func (s *pgStore) GetProductImagesByID(ctx context.Context, Id int) (pi []ProductImage, err error) {
	err = s.db.Select(&pi, getProductImagesByIDQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error Fetching ProductImages")
		return
	}
	return
}
