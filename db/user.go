package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"

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

	}()

	var dbUser User

	err = s.db.Get(&dbUser, getUserQuery, userID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("User Not found ")
		return
	}
	colName, colValue := PrepareParameters(dbUser, userProfile)
	var updateQuery string
	if len(colName) > 1 {
		updateQuery = `UPDATE users SET (` + strings.Join(colName, ",") + `)=('` + strings.Join(colValue, "','") + "') where id=$1"

	}
	if len(colName) == 1 {
		updateQuery = `UPDATE users SET ` + colName[0] + `='` + colValue[0] + "' where id=$1"
	}

	_, err = tx.ExecContext(ctx,
		updateQuery,

		userID,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating user profile")
		return
	}
	tx.Commit()
	updatedUser, err = s.GetUser(ctx, userID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting user from database with userID: ", userID)
		return
	}

	return

}

func PrepareParameters(userDb User, userProfile User) (colNames []string, colValues []string) {
	user := User{}
	elem := reflect.ValueOf(&userProfile).Elem()
	returnElem := reflect.ValueOf(&user).Elem()
	DbElem := reflect.ValueOf(&userDb).Elem()
	values := make([]interface{}, 0)
	for i := 0; i < elem.NumField(); i++ {
		if elem.Field(i).Interface() != returnElem.Field(i).Interface() {
			varname := DbElem.Type().Field(i).Tag.Get("db")
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
