package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
// @Description Show User User
// @Router /user/{id} [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var idParam = mux.Vars(req)["id"]
		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Error("Conversion Failed")
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
