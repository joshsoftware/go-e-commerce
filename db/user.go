package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/crypto/bcrypt"

	logger "github.com/sirupsen/logrus"
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

const (
	getUserQuery = `SELECT * from users where id=$1`
)

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

func (s *pgStore) UpdateUser(ctx context.Context, user User, userID int) (err error) {

	colName, colValue := PrepareParameters(user)
	var updateQuery string

	if len(colName) > 1 {
		updateQuery = fmt.Sprintf("Update users SET (%s)=('%s') where id =$1", strings.Join(colName, ","), strings.Join(colValue, "','"))
	}
	if len(colName) == 1 {
		updateQuery = fmt.Sprintf("Update users SET %s = '%s' where id =$1", colName[0], colValue[0])
	}

	_, err = s.db.Exec(
		updateQuery,
		userID,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error updating user profile")
		return
	}

	return

}

func (s *pgStore) AuthenticateUser(ctx context.Context, u User) (user User, err error) {

	err = s.db.Get(&user, "SELECT * FROM users where email = $1", u.Email)
	if err != nil {
		logger.WithField("err", err.Error()).Error("no such user available")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		logger.WithField("Error", err.Error())
	}
	return
}

func PrepareParameters(userProfile User) (colNames []string, colValues []string) {
	user := User{}
	elem := reflect.ValueOf(&userProfile).Elem()
	returnElem := reflect.ValueOf(&user).Elem()
	values := make([]interface{}, 0)
	for i := 0; i < elem.NumField(); i++ {
		if elem.Field(i).Interface() != returnElem.Field(i).Interface() {
			varname := returnElem.Type().Field(i).Tag.Get("db")
			colNames = append(colNames, varname)
			varValue := elem.Field(i).Interface()
			values = append(values, varValue)

		}

	}

	colValues = make([]string, len(values))
	for i, v := range values {
		colValues[i] = fmt.Sprintf("%s", v)
	}
	return

}
