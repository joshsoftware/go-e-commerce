package db

import (
	"context"
	"database/sql"

	logger "github.com/sirupsen/logrus"
)

const (
	insertUserQuery = `INSERT INTO users (first_name, last_name, email, mobile, country, state, city, address, password) 
	VALUES (:first_name, :last_name, :email, :mobile, :country, :state, :city, :address, :password)`

	getUserByEmailQuery = `SELECT * FROM users WHERE email=$1 LIMIT 1`
)

// User - struct representing a user
type User struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Mobile    string `db:"mobile" json:"mobile"`
	Country   string `db:"country" json:"country"`
	State     string `db:"state" json:"state"`
	City      string `db:"city" json:"city"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"password"`
	CreatedAt string `db:"created_at" json:"created_at"`
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
func (s *pgStore) CreateNewUser(ctx context.Context, u User) (newUser User, err error) {
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
	_, newUser, err = s.CheckUserByEmail(ctx, u.Email)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting user from database with email: " + u.Email)
		return
	}
	return
}

func (s *pgStore) CheckUserByEmail(ctx context.Context, email string) (check bool, user User, err error) {
	err = s.db.Get(&user, getUserByEmailQuery, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, user, err
		}
		logger.WithField("err", err.Error()).Error("Error while selecting user from database by email" + email)
		return
	}
	return true, user, err
}
