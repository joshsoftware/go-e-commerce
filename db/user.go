package db

import (
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
)

const (
	deleteUserQuery = `DELETE FROM users WHERE id=$1`
	getUserQuery    = `SELECT * from users where id=$1`
)

//User Struct for declaring attributes of User
type User struct {
	ID         int       `db:"id" json:"id"`
	FirstName  string    `db:"first_name" json:"first_name"`
	LastName   string    `db:"last_name" json:"last_name"`
	Email      string    `db:"email" json:"email"`
	Mobile     string    `db:"mobile" json:"mobile"`
	Address    string    `db:"address" json:"address"`
	Password   string    `db:"password" json:"password"`
	Country    string    `db:"country" json:"country"`
	State      string    `db:"state" json:"state"`
	City       string    `db:"city" json:"city"`
	IsAdmin    bool      `db:"isadmin" json:"isAdmin"`
	IsDisabled bool      `db:"isdisabled" json:"isDisabled"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func (s *pgStore) ListUsers(ctx context.Context) (users []User, err error) {
	err = s.db.Select(&users, "SELECT * FROM users ORDER BY name ASC")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing users")
		return
	}

	return
}

func (s *pgStore) GetUser(ctx context.Context, id int) (user User, err error) {
	err = s.db.Get(&user, getUserQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error(fmt.Errorf("error selecting user from database by id %d", id))
		return
	}
	return
}

func (s *pgStore) DeleteUserByID(ctx context.Context, userID int) (err error) {
	_, err = s.db.Exec(deleteUserQuery, userID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error Deleting User")
		return
	}
	return
}
