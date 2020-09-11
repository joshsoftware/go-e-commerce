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

// @Title UpdateUser
// @Description update User by id
// @Router /users/id [put]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func updateUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header["Token"]
		userID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		var user db.User

		err = json.NewDecoder(req.Body).Decode(&user)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error while decoding user")
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid json body",
				},
			})
			return
		}

		if user.Email != "" {

			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", "CAnnot update Email")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "cannot update email id !!",
				},
			})
			return
		}

		var updatedUser db.User
		updatedUser, err = deps.Store.UpdateUser(req.Context(), user, int(userID))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			logger.WithField("err", err.Error()).Error("Error while updating user's profile")
			return
		}
		responses(rw, http.StatusOK, successResponse{Data: updatedUser})

		return

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
		check, _, err := deps.Store.CheckUserByEmail(req.Context(), user.Email)

		// If check true then user is already registered
		if check {
			e := errorResponse{
				Error: "user already registered",
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
		_, err = deps.Store.CreateNewUser(req.Context(), user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in inserting user in database")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := successResponse{
			Data: "user successfully registered",
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

// @Title getUser
// @Description get User by id
// @Router /users/{id} [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//fetch usedId from request
		authToken := req.Header["Token"]
		userID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		user, err1 := deps.Store.GetUser(req.Context(), int(userID))
		if err1 != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		respBytes, err := json.Marshal(user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error marshaling users data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}
