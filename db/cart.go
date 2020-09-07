package db

import(
	"context"
	logger "github.com/sirupsen/logrus"
)

func (s *pgStore) AddToCart(ctx context.Context, cartId, productId int) (err error) {
	insert := `INSERT INTO cart (id, product_id, quantity) VALUES ($1, $2, 1)`
	_, err = s.db.Exec(insert, cartId, productId)
	if err != nil{
		logger.WithField("err", err.Error()).Error("Error adding to cart")
		return
	}
	return
}

func (s *pgStore) RemoveFromCart(ctx context.Context, cartId, productId int) (err error) {
	delete := `DELETE FROM cart WHERE id = $1 AND product_id = $2`
	_, err = s.db.Exec(delete, cartId, productId)
	if err != nil{
		logger.WithField("err", err.Error()).Error("Error while removing from cart")
		return
	}
	return
}

func (s *pgStore) UpdateIntoCart(ctx context.Context, quantity, cartId, productId int) (err error) {
	update := `UPDATE cart SET quantity = $1 WHERE id = $2 AND product_id = $3`
	_, err = s.db.Exec(update, quantity, cartId, productId)
	if err != nil{
		logger.WithField("err", err.Error()).Error("Error while updating into cart")
		return
	}
	return
}