package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"joshsoftware/go-e-commerce/config"
	"net/http"
)

const (
	versionHeader = "Accept"
)

/*InitRouter is  The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter(deps Dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	//Route for User Login
	router.HandleFunc("/login", userLoginHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	//Router for Get User from ID
	router.Handle("/user", jwtMiddleWare(getUserHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	//Router for User Logout
	router.Handle("/logout", jwtMiddleWare(userLogoutHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	//Router for Get All Users
	router.HandleFunc("/users", listUsersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	//routes for cart operations
	router.Handle("/cart", jwtMiddleWare(addToCartHandler(deps), deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(deleteFromCartHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(updateIntoCartHandler(deps), deps)).Methods(http.MethodPut).Headers(versionHeader, v1)
	return
}

//jwtMiddleWare function is used to authenticate and authorize the incoming request
func jwtMiddleWare(endpoint http.Handler, deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header.Get("Token")

		//Checking if token not present in header
		if len(authToken) < 1 {
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Missing Authorization Token",
				},
			})
			return
		}

		_, _, err := getDataFromToken(authToken)
		if err != nil {
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Unauthorized User",
				},
			})
			return
		}

		//Fetching Status of Token Being Blacklisted or Not
		// Unauthorized User if Token BlackListed
		if isBlacklisted, _ := deps.Store.CheckBlacklistedToken(req.Context(), authToken); isBlacklisted {
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Unauthorized User",
				},
			})
			return
		}
		endpoint.ServeHTTP(rw, req)
	})
}
