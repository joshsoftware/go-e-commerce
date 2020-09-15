package db

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	logger "github.com/sirupsen/logrus"
)

const (
	updateUserQuery = `UPDATE users SET (
		first_name,
		last_name,
		mobile,
		address,
		password,
		country,
		state,
		city
		
		) = 
		($1, $2, $3, $4, $5, $6 ,$7,$8) where id = $9 `

	getUserQuery = `SELECT * from users where id=$1`
)

//User is a structure of the user
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
	err = s.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		logger.WithField("err", err.Error()).Error("error listing users")
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

func (s *pgStore) UpdateUserByID(ctx context.Context, user User, userID int) (err error) {

	//check if password is to be updated then convert it to hashcode
	if user.Password != "" {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while creating hash of the password")
			return err
		}
		user.Password = string(hashedPassword)

	}

	_, err = s.db.Exec(
		updateUserQuery,
		user.FirstName,
		user.LastName,
		user.Mobile,
		user.Address,
		user.Password,
		user.Country,
		user.State,
		user.City,
		userID,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating user profile")
		return
	}

	return

}

//Validate function for user
func (user *User) Validate() (valid bool) {
	fieldErrors := make(map[string]string)

	if user.FirstName == "" {
		fieldErrors["name"] = "Can't be blank"
	}
	if user.LastName == "" {
		fieldErrors["LastName"] = "Can't be blank"
	}
	if user.Mobile == "" {
		fieldErrors["Mobile"] = "Can't be blank"
	}
	if user.Password == "" {
		fieldErrors["Password"] = "Can't be blank"
	}
	if user.Address == "" {
		fieldErrors["Address"] = "Can't be blank"
	}
	if user.Country == "" {
		fieldErrors["Country"] = "Can't be blank"
	}
	if user.State == "" {
		fieldErrors["State"] = "Can't be blank"
	}
	if user.City == "" {
		fieldErrors["City"] = "Can't be blank"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	valid = false
	//TODO: Ask what other validations are expected

	return
}

//TODO add function for aunthenticating user
