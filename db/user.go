package db

import (
	"context"
	"fmt"

	logger "github.com/sirupsen/logrus"
)

//User is a structure of the user
type User struct {
	UserID       int    `db:"userid" json:"id"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	Email        string `db:"email" json:"email"`
	MobileNumber string `db:"mobile_number" json:"mobile_number"`
	Password     string `db:"password" json:"password"`
	Country      string `db:"country" json:"country"`
	State        string `db:"state" json:"state"`
	City         string `db:"city" json:"city"`
	Address      string `db:"address" json:"address"`
}

const (
	updateUserQuery = `UPDATE users SET (
	first_name,
	last_name,
	email,
	mobile_number,
	password,
	country,
	state,
	city,
	address
	) = 
	($1, $2, $3, $4, $5,$6,$7,$8,$9) where userid = $10 `

	getUserQuery = `SELECT * from users where userid=$1`
)

func (s *pgStore) ListUsers(ctx context.Context) (users []User, err error) {
	err = s.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing users")
		return
	}

	return
}

func (s *pgStore) GetUser(ctx context.Context, id int) (user User, err error) {
	err = s.db.Get(&user, getUserQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting user from database by id " + fmt.Sprint(id))
		return
	}

	return
}

func (s *pgStore) UpdateUser(ctx context.Context, userProfile User, userID int) (updatedUser User, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating transaction")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	var dbUser User

	err = s.db.Get(&dbUser, getUserQuery, userID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("User Not found ")
		return
	}

	_, err = s.db.Exec(
		updateUserQuery,
		userProfile.FirstName,
		userProfile.LastName,
		userProfile.Email,
		userProfile.MobileNumber,
		userProfile.Password,
		userProfile.Country,
		userProfile.State,
		userProfile.City,
		userProfile.Address,
		userID,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating user profile")
		return
	}

	updatedUser, err = s.GetUser(ctx, userID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting user from database with userID: ", userID)
		return
	}

	return

}
