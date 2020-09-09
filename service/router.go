package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"joshsoftware/go-e-commerce/config"
	"net/http"
	"strconv"
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
	router.Handle("/user/{id}", jwtMiddleWare(getUserHandler(deps), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	//Router for User Logout
	router.Handle("/user/{id}/logout", jwtMiddleWare(userLogoutHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	//Router for Get All Users
	router.HandleFunc("/users", listUsersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	return
}

//jwtMiddleWare function is used to authenticate and authorize the incoming request
func jwtMiddleWare(endpoint http.Handler, deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		mySigningKey := config.JWTKey()

		//Fetching userID from RequestURL
		var idParam = mux.Vars(req)["id"]
		validID, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Error(err.Error())
		}

		authToken := req.Header["Token"]

		//Checking if token not present in header
		if authToken == nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		token, err := jwt.Parse(authToken[0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		//Fetching Status of Token Being Blacklisted or Not
		// Unauthorized User if Token BlackListed
		if isBlacklisted, _ := deps.Store.CheckBlacklistedToken(req.Context(), authToken[0]); isBlacklisted {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		//Unauthorized User if Token Invalid
		if !ok && !token.Valid {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		userID := claims["id"]

		//Unauthorized User if userID in Token Doesn't Match userID in RequestURL
		if float64(validID) != userID {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
			return
		}

		endpoint.ServeHTTP(rw, req)
	})
}
