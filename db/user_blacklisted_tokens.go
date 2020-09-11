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
	UserID         int       `db:"user_id" json:"user_id"`
	Token          string    `db:"token" json:"token"`
	ExpirationDate time.Time `db:"expiration_date" json:"expiration_date"`
}

const (
	insertBlacklistedToken = `INSERT INTO user_blacklisted_tokens
(user_id, token, expiration_date)
VALUES ($1, $2, $3)`
)

//CreateBlacklistedToken function to insert the blacklisted token in database
func (s *pgStore) CreateBlacklistedToken(ctx context.Context, token BlacklistedToken) (err error) {
	_, err = s.db.Exec(insertBlacklistedToken, token.UserID, token.Token, token.ExpirationDate)

	if err != nil {
		errMsg := fmt.Sprintf("Error inserting the blacklisted token for user with id %v", token.UserID)
		logger.WithField("err", err.Error()).Error(errMsg)
		return
	}
	return
}

//CheckBlacklistedToken function to check if token is blacklisted earlier
func (s *pgStore) CheckBlacklistedToken(ctx context.Context, token string) (bool, int) {

	var userID int
	query1 := fmt.Sprintf("SELECT user_id FROM user_blacklisted_tokens WHERE token='%s'", token)
	err := s.db.QueryRow(query1).Scan(&userID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Either Query Failed or No Rows Found")
		return false, -1
	}
	return true, userID
}
