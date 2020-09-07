package db

import (
	"context"
	"database/sql"
	logger "github.com/sirupsen/logrus"
)

const (
	getUserByMobileQuery = `SELECT * FROM users WHERE email=$1 LIMIT 1`
)

type User struct {
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Mobile    string `db:"mobile" json:"mobile"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"password"`
}

func (s *pgStore) ListUsers(ctx context.Context) (users []User, err error) {
	err = s.db.Select(&users, "SELECT * FROM users ORDER BY first_name ASC")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing users")
		return
	}

	return
}

// GetUserByMobile - Given a mobile number, return that user.
func (s *pgStore) GetUserByMobile(ctx context.Context, mobile string) (user User, err error) {
	err = s.db.Get(&user, getUserByMobileQuery, mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("User with mobile number " + mobile + "doenot Exist")
		}
		// Possible that there's no rows in the result set
		logger.WithField("err", err.Error()).Error("Error selecting user from database by mobile number " + mobile)
		return
	}
	return
}
