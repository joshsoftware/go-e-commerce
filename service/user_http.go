package service

import (
	"encoding/json"
	"joshsoftware/go-e-commerce/db"
	"net/http"

	logger "github.com/sirupsen/logrus"
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

// @Title getUser
// @Description get User by id
// @Router /users/{id} [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header["Token"]
		userID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		user, err := deps.Store.GetUser(req.Context(), int(userID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			repsonse(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "Id Not Found",
				},
			})
			return
		}
		repsonse(rw, http.StatusOK, successResponse{Data: user})
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
			repsonse(rw, http.StatusBadRequest, errorResponse{
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
			repsonse(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			logger.WithField("err", err.Error()).Error("Error while updating user's profile")
			return
		}
		repsonse(rw, http.StatusOK, successResponse{Data: updatedUser})

		return

	})

}
