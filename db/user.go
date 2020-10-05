package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"time"

	"golang.org/x/crypto/bcrypt"

	logger "github.com/sirupsen/logrus"
)

const (
	insertUserQuery = `INSERT INTO users (first_name, last_name, email, mobile, country, state, city, address, password) 
	VALUES (:first_name, :last_name, :email, :mobile, :country, :state, :city, :address, :password)`

	getUserByEmailQuery = `SELECT * FROM users WHERE email=$1 LIMIT 1`
	updateUserQuery     = `UPDATE users SET (
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

	deleteUserQuery       = `DELETE FROM users WHERE id=$1`
	disableUserQuery      = `UPDATE users SET isdisabled =$1 WHERE id=$2`
	enableUserQuery       = `UPDATE users SET isdisabled =$1 WHERE id=$2`
	deleteUsersTokenQuery = `DELETE FROM user_blacklisted_tokens WHERE user_id=$1`
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
	IsVerified bool      `db:"isverified" json:"isVerified"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

//ListUsers function to fetch all Users From Database
func (s *pgStore) ListUsers(ctx context.Context) (users []User, err error) {
	err = s.db.Select(&users, "SELECT * FROM users ORDER BY id ASC")
	if err != nil {
		logger.WithField("err", err.Error()).Error(fmt.Errorf("error selecting users from database"))
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

//AuthenticateUser Function checks if User has Registered before Login
// and Has Entered Correct Credentials
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

	// creating hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userProfile.Password), 8)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while creating hash of the password")
		return
	}
	userProfile.Password = string(hashedPassword)

	var dbUser User

	err = s.db.Get(&dbUser, getUserQuery, userID)

	if err != nil {
		logger.WithField("err", err.Error()).Error("User Not found ")
		return
	}
	colName, colValue := prepareParameters(dbUser, userProfile)
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

func prepareParameters(userDb User, userProfile User) (colNames []string, colValues []string) {
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
func (s *pgStore) UpdateUserByID(ctx context.Context, user User, userID int) (err error) {

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

//ValidatePatchParams function for user
func (user *User) ValidatePatchParams(u User) (err error) {

	if u.FirstName != "" {
		user.FirstName = u.FirstName
	}
	if u.LastName != "" {
		user.LastName = u.LastName
	}
	if u.Mobile != "" {
		user.Mobile = u.Mobile
	}
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while creating hash of the password")
			return err
		}
		user.Password = string(hashedPassword)
	}
	if u.Address != "" {
		user.Address = u.Address
	}
	if u.Country != "" {
		user.Country = u.Country
	}
	if u.State != "" {
		user.State = u.State
	}
	if u.City != "" {
		user.City = u.City
	}
	return
}

func (s *pgStore) DisableUserByID(ctx context.Context, userID int) (err error) {
	_, err = s.db.Exec(disableUserQuery, true, userID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error disabling User")
		return
	}
	return
}

func (s *pgStore) EnableUserByID(ctx context.Context, userID int) (err error) {
	_, err = s.db.Exec(enableUserQuery, false, userID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error enabling User")
		return
	}
	return
}

func (s *pgStore) DeleteUserByID(ctx context.Context, userID int) (err error) {
	_, err = s.db.Exec(deleteUsersTokenQuery, userID)
	_, err = s.db.Exec(deleteUserQuery, userID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error Deleting User")
		return
	}
	return
}
