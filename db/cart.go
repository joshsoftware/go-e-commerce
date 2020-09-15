package db

import(
	"context"
	logger "github.com/sirupsen/logrus"
)

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