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
			logger.WithField("err", err.Error()).Error("error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(users)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error marshaling users data")
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
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			repsonse(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
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
			logger.WithField("err", err.Error()).Error("error while decoding user")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "invalid json body",
				},
			})
			return
		}

		if user.Email != "" {

			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", "Cannot update email")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "cannot update email id !!",
				},
			})
			return
		}

		err = deps.Store.UpdateUser(req.Context(), user, int(userID))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			repsonse(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "internal server error",
				},
			})
			logger.WithField("err", err.Error()).Error("error while updating user's profile")
			return
		}

		updatedUser, err := deps.Store.GetUser(req.Context(), int(userID))
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			repsonse(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "error while fetching users",
				},
			})
			return
		}
		repsonse(rw, http.StatusOK, successResponse{Data: updatedUser})

		return

	})

}
