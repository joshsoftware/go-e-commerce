package db

import (
	"context"
	"database/sql"

	logger "github.com/sirupsen/logrus"
)

const (
	insertUserQuery = `INSERT INTO users (first_name, last_name, email, password, mobile_number, country, state, city, address) 
	VALUES (:first_name, :last_name, :email, :password, :mobile_number, :country, :state, :city, :address)`

	getUserByEmailQuery = `SELECT * FROM users WHERE email=$1 LIMIT 1`
)

// User - struct representing a user
type User struct {
	UserID       int    `db:"userid" json:"user_id"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	Email        string `db:"email" json:"email"`
	Password     string `db:"password" json:"password"`
	MobileNumber string `db:"mobile_number" json:"mobile_number"`
	Country      string `db:"country" json:"country"`
	State        string `db:"state" json:"state"`
	City         string `db:"city" json:"city"`
	Address      string `db:"address" json:"address"`
}

func (s *pgStore) ListUsers(ctx context.Context) (users []User, err error) {
	err = s.db.Select(&users, "SELECT * FROM users ORDER BY name ASC")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing users")
		return
	}

	return
}

// CreateNewUser = creates a new user in database
func (s *pgStore) CreateNewUser(ctx context.Context, u User) (err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in beginning user insert transaction")
		return
	}
	_, err = tx.NamedExec(insertUserQuery, u)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while inserting user into database")
		return
	}
	err = tx.Commit()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while commiting transaction inserting user")
		return
	}
	return
}

func (s *pgStore) CheckUserByEmail(ctx context.Context, email string) (check bool, err error) {
	user := User{}
	err = s.db.Get(&user, getUserByEmailQuery, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		logger.WithField("err", err.Error()).Error("Error while selecting user from database by email" + email)
		return
	}
	return true, err
}
