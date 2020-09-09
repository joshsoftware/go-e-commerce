package service

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"joshsoftware/go-e-commerce/db"
	"net/http"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type errorResponse struct {
	Error string `json:"error"`
}
type successResponse struct {
	Message string `json:"message"`
}

// @Title listUsers
// @Description list all User
// @Router /users [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func listUsersHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		users, err := deps.Store.ListUsers(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(users)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error marshaling users data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

// @Title registerUser
// @Description registers new user
// @Router /register [post]
// @Accept  json
// @Success 201 {object}
// @Failure 400 {object}
func registerUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// reading data from body
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in reading request body")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		user := db.User{}
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while Unmarshalling request json")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// For checking if user already registered
		check, err := deps.Store.CheckUserByEmail(req.Context(), user.Email, user.Mobile)

		// If check true then user is already registered
		if check {
			e := errorResponse{
				Error: "user's email or mobile already registered",
			}
			respBytes, err := json.Marshal(e)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error while marshalling error msg ")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(respBytes)
			return
		}

		// For checking error occured while looking already registered user
		if err != nil && err != sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("Error while looking existing user")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// creating hash of the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while creating hash of the password")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		// Storing new user's data in database
		err = deps.Store.CreateNewUser(req.Context(), user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in inserting user in database")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := successResponse{
			Message: "user successfully registered",
		}
		respBytes, err := json.Marshal(msg)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while marshalling success msg ")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(respBytes)
		return

	})
}
