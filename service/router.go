package service

import (
	"context"
	"fmt"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
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

	// Route for register new user
	router.HandleFunc("/register", registerUserHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	//Route for google Oauth
	router.HandleFunc("/auth/google", handleAuth(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	//Route for User Login
	router.HandleFunc("/login", userLoginHandler(deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	//Router for User Logout
	router.Handle("/logout", jwtMiddleWare(userLogoutHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	//Router for users operations
	router.Handle("/user", jwtMiddleWare(userMiddleware(getUserHandler(deps), deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/admin", jwtMiddleWare(adminMiddleware(getUserHandler(deps), deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/users", jwtMiddleWare(adminMiddleware(listUsersHandler(deps), deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/user/update", jwtMiddleWare(userMiddleware(updateUserHandler(deps), deps), deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)
	router.Handle("/admin/update", jwtMiddleWare(adminMiddleware(updateUserHandler(deps), deps), deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)
	router.Handle("/user/{id:[0-9]+}", jwtMiddleWare(adminMiddleware(deleteUserHandler(deps), deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.Handle("/user/disable/{id:[0-9]+}", jwtMiddleWare(adminMiddleware(disableUserHandler(deps), deps), deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)
	router.Handle("/user/enable/{id:[0-9]+}", jwtMiddleWare(adminMiddleware(enableUserHandler(deps), deps), deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)

	// routes for product operations
	router.HandleFunc("/products", listProductsHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/product/{product_id:[0-9]+}", getProductByIdHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/createProduct", jwtMiddleWare(adminMiddleware(createProductHandler(deps), deps), deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.Handle("/product/{product_id:[0-9]+}", jwtMiddleWare(adminMiddleware(deleteProductByIdHandler(deps), deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.HandleFunc("/products/filters", getProductByFiltersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/product/stock", updateProductStockByIdHandler(deps)).Methods(http.MethodPut).Headers(versionHeader, v1)
	router.HandleFunc("/products/search", getProductBySearchHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/product/{product_id:[0-9]+}", updateProductByIdHandler(deps)).Methods(http.MethodPut).Headers(versionHeader, v1)
	router.PathPrefix("/static/products").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/"))))

	//routes for cart operations
	router.Handle("/cart", jwtMiddleWare(userMiddleware(getCartHandler(deps), deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(userMiddleware(addToCartHandler(deps), deps), deps)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(userMiddleware(removeFromCartHandler(deps), deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.Handle("/cart", jwtMiddleWare(userMiddleware(updateIntoCartHandler(deps), deps), deps)).Methods(http.MethodPut).Headers(versionHeader, v1)

	router.HandleFunc("/footer", getFooterHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/country_data", countryDataHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))

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

func getUserFromToken(ctx context.Context, deps Dependencies, authToken string) (user db.User, err error) {

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

		user, err := getUserFromToken(req.Context(), deps, authToken)

		if user.IsDisabled {
			responses(rw, http.StatusForbidden, errorResponse{
				Error: messageObject{
					Message: "User Forbidden From Accesing Data",
				},
			})
			return
		}

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
		user, err := getUserFromToken(req.Context(), deps, authToken)

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
