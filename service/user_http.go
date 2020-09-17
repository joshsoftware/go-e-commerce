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

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
				},
			})
			return
		}

		//check if the users is admin
		if user.IsAdmin == true {
			logger.WithField("err", "cannot delete an admin")
			rw.WriteHeader(http.StatusForbidden)
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "access denied for delete",
				},
			})
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

// @Title disableuser
// @Description disable the user record by given id
// @Router /user/disable/{id} [Patch]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func disableUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		//TODO get the IsAdmin field of the user from the token
		isAdmin := true
		if !isAdmin {
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "User Forbidden From disabling Data",
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

		//check if the provided id is not of another Admin
		user, err := deps.Store.GetUser(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error while fetching User")
			rw.WriteHeader(http.StatusNotFound)
			responses(rw, http.StatusNotFound, errorResponse{
				Error: messageObject{
					Message: "id Not Found",
				},
			})
			return
		}

		//check if the user is admin
		if user.IsAdmin == true {
			logger.WithField("err", "cannot disable an admin")
			rw.WriteHeader(http.StatusForbidden)
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "access denied for disable",
				},
			})
			return
		}

		err = deps.Store.DisableUserByID(req.Context(), userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("error in disabling user data")
			rw.WriteHeader(http.StatusBadRequest)
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "error in disabling user data",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: "record disabled Successfully",
		})

		return

	})
}
