package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

//listUsersHandler function fetch all users from database
// and return as json object
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

//listUsersHandler function fetch specific user from database
// and return as json object
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//fetch usedId from request
		var idParam = mux.Vars(req)["id"]

		id, err := strconv.Atoi(idParam)

		if err != nil {
			logger.WithField("err", err.Error()).Error("Invalid User ID")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := deps.Store.GetUser(req.Context(), id)
		if err != nil {
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
