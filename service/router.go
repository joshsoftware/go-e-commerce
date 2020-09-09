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

	//Router for Get User from ID
	router.Handle("/user/{id}/logout", jwtMiddleWare(userLogoutHandler(deps), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	//Router for Get All Users
	router.HandleFunc("/users", listUsersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	return
}

func jwtMiddleWare(endpoint http.Handler, deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		mySigningKey := config.JWTKey()

		var idParam = mux.Vars(req)["id"]
		validID, ok := strconv.Atoi(idParam)

		if ok != nil {
			logger.Error(ok.Error())
		}
		authToken := req.Header["Token"]
		if authToken != nil {
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

			disdis := deps.Store.CheckBlacklistedToken(req.Context(), authToken[0])

			fmt.Println("DIsDis : ", disdis)
			if !ok && !token.Valid && !disdis {
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Write([]byte("Unauthorized"))
				return
			}

			userID := claims["id"]

			if float64(validID) != userID {
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Write([]byte("Unauthorized"))
				return
			}

			endpoint.ServeHTTP(rw, req)
		}
	})
}
