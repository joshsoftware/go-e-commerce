package service

import (
	"fmt"
	"joshsoftware/go-e-commerce/config"
	"net/http"

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

	//Route for User Login
	router.HandleFunc("/login", userLoginHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	//Router for User Logout
	router.Handle("/logout", jwtMiddleWare(userLogoutHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	//Router for Get All Users
	router.Handle("/users", jwtMiddleWare(listUsersHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	//Router for Get User from ID
	router.Handle("/user", jwtMiddleWare(getUserHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/user/update", jwtMiddleWare(updateUserHandler(deps), deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)

	//Route for google Oauth
	router.HandleFunc("/auth/google", handleAuth(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/register", registerUserHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/footer", getFooterHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/products", listProductsHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/product/{product_id:[0-9]+}", getProductByIdHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/products/category/{category_id:[0-9]+}", listProductsByCategoryHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/createProduct", createProductHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/product/{product_id:[0-9]+}", deleteProductByIdHandler(deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.Handle("/cart", jwtMiddleWare(getCartHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.HandleFunc("/users", listUsersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	//routes for cart operations
	router.Handle("/cart", jwtMiddleWare(addToCartHandler(deps), deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(removeFromCartHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(updateIntoCartHandler(deps), deps)).Methods(http.MethodPut).Headers(versionHeader, v1)
	return
}

//jwtMiddleWare function is used to authenticate and authorize the incoming request
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
