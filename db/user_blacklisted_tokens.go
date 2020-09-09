package db

import (
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
)

//BlacklistedToken - struct representing a token to be blacklisted (logout)
type BlacklistedToken struct {
	ID             int       `db:"id" json:"id"`
	UserID         float64   `db:"user_id" json:"user_id"`
	Token          string    `db:"token" json:"token"`
	ExpirationDate time.Time `db:"expiration_date" json:"expiration_date"`
}

const (
	insertBlacklistedToken = `INSERT INTO user_blacklisted_tokens
(user_id, token, expiration_date)
VALUES ($1, $2, $3)`

	selectBlacklistedToken = `SELECT * FROM user_blacklisted_tokens
WHERE token=$1`
)

func (s *pgStore) CreateBlacklistedToken(ctx context.Context, token BlacklistedToken) (err error) {
	_, err = s.db.Exec(insertBlacklistedToken, token.UserID, token.Token, token.ExpirationDate)

	if err != nil {
		errMsg := fmt.Sprintf("Error inserting the blacklisted token for user with id %d", token.UserID)
		logger.WithField("err", err.Error()).Error(errMsg)
		return
	}
	return
}

func (s *pgStore) CheckBlacklistedToken(ctx context.Context, token string) bool {
	query := fmt.Sprintf("SELECT * FROM user_blacklisted_tokens WHERE token=%s", token)
	result, err := s.db.Exec(query)
	fmt.Println(result)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Query Failed")
		return false
	}
	return true
}
