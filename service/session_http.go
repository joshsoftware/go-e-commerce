package service

import (
	"encoding/json"
	"fmt"
	ae "joshsoftware/go-e-commerce/apperrors"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
)

//AuthBody stores responce body for login
type authBody struct {
	Message string `json:"meassage"`
	Token   string `json:"token"`
}

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
			responses(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "JSON Decoding Failed",
				},
			})

			return
		}

		//TODO change no need to return user object from Authentication
		//checking if the user is authenticated or not
		// by passing the credentials to the AuthenticateUser function
		user, err1 := deps.Store.AuthenticateUser(req.Context(), user)
		if err1 != nil {
			logger.WithField("err", err1.Error()).Error("Invalid Credentials")
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Invalid Credentials",
				},
			})
			return
		}

		//Generate new JWT token if the user is authenticated
		// and return the token in request header
		token, err := generateJwt(user.ID)
		if err != nil {
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Token Generation Failure",
				},
			})
			return
		}

		responses(rw, http.StatusOK, successResponse{
			Data: authBody{
				Message: "Login Successfull",
				Token:   token,
			},
		})
	})
}

//userLogoutHandler function logs the user off
// and add the valid JWT token in BlacklistedToken
func userLogoutHandler(deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		//fetching the token from header
		authToken := req.Header["Token"]

		//fetching details from the token
		userID, expirationTimeStamp, err := getDataFromToken(authToken[0])
		if err != nil {
			responses(rw, http.StatusUnauthorized, errorResponse{
				Error: messageObject{
					Message: "Unauthorized User",
				},
			})
			return
		}
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
			responses(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal Server Error",
				},
			})
			return
		}
		responses(rw, http.StatusOK, successResponse{
			Data: messageObject{
				Message: "Logged Out Successfully",
			},
		})
		return
	})
}

func getDataFromToken(Token string) (userID float64, expirationTime int64, err error) {
	mySigningKey := config.JWTKey()

	//Checking if token not present in header
	if len(Token) < 1 {
		ae.Error(ae.ErrMissingAuthHeader, "Missing Authentication Token From Header", err)
		return
	}

	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})
	if err != nil {
		ae.Error(ae.ErrInvalidToken, "Invalid Token", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		ae.Error(ae.ErrInvalidToken, "Invalid Token", err)
		return
	}

	userID = claims["id"].(float64)
	expirationTime = int64(claims["exp"].(float64))
	return
}
