package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

// @Title deleteuser
// @Description delete the user record by given id
// @Router /user/{id} [Delete]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func deleteUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		//TODO get the IsAdmin field of the user from the token
		isAdmin := true
		if !isAdmin {
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "User Forbidden From Deleting Data",
				},
			})
			return
		}

		params := mux.Vars(req)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in getting id")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		err = deps.Store.DeleteUserByID(req.Context(), userID)

		if err != nil {
			logger.WithField("err", err.Error()).Error("error in deleting user data")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "error in deleting user data",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "record deleted Successfully",
		})

		return

	})
}
