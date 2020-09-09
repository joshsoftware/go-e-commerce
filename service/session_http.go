package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
	ae "joshsoftware/go-e-commerce/apperrors"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"time"
)

//generateJWT function generates and return a new JWT token
func generateJwt(userID int) (tokenString string, err error) {
	mySigningKey := config.JWTKey()
	if mySigningKey == nil {
		ae.Error(ae.ErrNoSigningKey, "Application error: No signing key configured", err)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Duration(config.JWTExpiryDurationHours()) * time.Hour).Unix()

	tokenString, err = token.SignedString(mySigningKey)
	if err != nil {
		ae.Error(ae.ErrSignedString, "Failed To Get Signed String", err)
		return
	}
	return
}

//userLoginHandler function take credentials in json
// and check if the credentials are correct
// also generate and returns a JWT token in the case of correct crendential
func userLoginHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		user := db.User{}

		//fetching the json object to get crendentials of users
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("JSON Decoding Failed")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		//checking if the user is authenticated or not
		// by passing the credentials to the AuthenticateUser function
		user, err1 := deps.Store.AuthenticateUser(req.Context(), user)
		if err1 != nil {
			logger.WithField("err", err1.Error()).Error("Invalid Credentials")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Generate new JWT token if the user is authenticated
		// and return the token in request header
		token, err := generateJwt(user.ID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Token Generation Failure"))
			return
		}

		respBytes, err := json.Marshal(user)
		if err != nil {
			ae.Error(ae.ErrJSONParseFail, "JSON Parsing Failed", err)
		}

		rw.Header().Add("Token", token)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

//userLogoutHandler function logs the user off
// and add the valid JWT token in BlacklistedToken
func userLogoutHandler(deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		mySigningKey := config.JWTKey()

		//fetching the token from header
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

		//Checking if token not valid
		if !ok && !token.Valid {
			ae.Error(ae.ErrInvalidToken, "Authentication Token Invalid", err)
			ae.JSONError(rw, http.StatusUnauthorized, err)
			return
		}

		//fetching details from the token
		userID := claims["id"].(float64)
		expirationTimeStamp := int64(claims["exp"].(float64))
		expirationDate := time.Unix(expirationTimeStamp, 0)

		//create a BlacklistedToken to add in database
		// To blacklist a user valid token
		userBlackListedToken := db.BlacklistedToken{
			UserID:         userID,
			ExpirationDate: expirationDate,
			Token:          authToken[0],
		}

		err = deps.Store.CreateBlacklistedToken(req.Context(), userBlackListedToken)
		if err != nil {
			ae.Error(ae.ErrFailedToCreate, "Error creating blaclisted token record", err)
			rw.Header().Add("Content-Type", "application/json")
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		return
	})
}
