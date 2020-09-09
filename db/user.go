package db

import (
	"context"
	"database/sql"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	ae "joshsoftware/go-e-commerce/apperrors"
)

//User Struct for declaring attributes of User
type User struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Mobile    string `db:"mobile" json:"mobile"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"password"`
	Country   string `db:"country" json:"country"`
	State     string `db:"state" json:"state"`
	City      string `db:"city" json:"city"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

func (s *pgStore) ListUsers(ctx context.Context) (users []User, err error) {
	err = s.db.Select(&users, "SELECT * FROM users ORDER BY first_name ASC")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing users")
		return
	}

	return
}

func (s *pgStore) GetUser(ctx context.Context, id int) (user User, err error) {

	err = s.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ae.ErrRecordNotFound
		}
		logger.WithField("err", err.Error()).Error("Query Failed")
		return
	}

	return
}

// GetUserByMobile - Given a mobile number, return that user.
func (s *pgStore) AuthenticateUser(ctx context.Context, u User) (user User, err error) {

	err = s.db.Get(&user, "SELECT * FROM users where email = $1", u.Email)
	if err != nil {
		logger.WithField("err", err.Error()).Error("No such User Available")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		logger.WithField("Error", err.Error())
	}

	return
}
