package service

import (
	"fmt"
	"net/http"

	"joshsoftware/go-e-commerce/config"

	"github.com/gorilla/mux"
)

const (
	versionHeader = "Accept"
)

//InitRouter :The routing mechanism. Mux helps us define handler functions and the access methods
func InitRouter(deps Dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	router.Handle("/users", jwtMiddleWare(listUsersHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/user", jwtMiddleWare(getUserHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/user/update", jwtMiddleWare(updateUserHandler(deps), deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)
	//Route for User Login
	router.HandleFunc("/login", userLoginHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	return
}

func jwtMiddleWare(endpoint http.Handler, deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header["Token"]

		_, _, err := getDataFromToken(authToken[0])
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		//Fetching Status of Token Being Blacklisted or Not
		// Unauthorized User if Token BlackListed
		if isBlacklisted, _ := deps.Store.CheckBlacklistedToken(req.Context(), authToken[0]); isBlacklisted {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}
		endpoint.ServeHTTP(rw, req)
	})
}
