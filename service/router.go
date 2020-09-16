package service

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
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
	router.Handle("/user", jwtMiddleWare(userMiddleware(getUserHandler(deps), deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	//Router for User Logout
	router.Handle("/logout", jwtMiddleWare(userLogoutHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	//Router for Get All Users
	router.Handle("/users", jwtMiddleWare(adminMiddleware(listUsersHandler(deps), deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
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

func helper(ctx context.Context, deps Dependencies, authToken string) (user db.User, err error) {

	userID, _, _, err := getDataFromToken(authToken)
	if err != nil {
		return user, err
	}

	user, err = deps.Store.GetUser(ctx, int(userID))
	if err != nil {
		return user, err
	}

	return user, nil
}

func userMiddleware(endpoint http.Handler, deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header.Get("Token")

		user, err := helper(req.Context(), deps, authToken)

		if user.IsAdmin || err != nil {
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

func adminMiddleware(endpoint http.Handler, deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header.Get("Token")
		user, err := helper(req.Context(), deps, authToken)

		if !user.IsAdmin || err != nil {
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
