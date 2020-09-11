package service

import (
	logger "github.com/sirupsen/logrus"
	"net/http"
)

//listUsersHandler function fetch all users from database
// and return as json object
func listUsersHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		users, err := deps.Store.ListUsers(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: users,
		})
	})
}

//listUsersHandler function fetch specific user from database
// and return as json object
func getUserHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//fetch usedId from request
		authToken := req.Header["Token"]
		userID, _, err := getDataFromToken(authToken[0])
		if err != nil {
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Unauthorized User",
				},
			})
			return
		}

		user, err1 := deps.Store.GetUser(req.Context(), int(userID))
		if err1 != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: user,
		})
	})
}
