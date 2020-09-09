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

// @Title userLogin
// @Description Logging User in
// @Router /login [post]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func userLoginHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		user := db.User{}
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("JSON Decoding Failed")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err1 := deps.Store.AuthenticateUser(req.Context(), user)
		if err1 != nil {
			logger.WithField("err", err1.Error()).Error("Invalid Credentials")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := generateJwt(user.ID)
		if err != nil {
			ae.Error(ae.ErrUnknown, "Unknown/unexpected error while creating JWT", err)
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}

		respBytes, err := json.Marshal(user)

		rw.Header().Add("Authorization", token)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

func userLogoutHandler(deps Dependencies) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		mySigningKey := config.JWTKey()

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

			if !ok && !token.Valid {
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Write([]byte("Unauthorized"))
				return
			}

			userID := claims["id"].(float64)

			if err != nil {
				logger.WithField("err", err.Error()).Error("Conversion Failed")
			}
			expirationTimeStamp := int64(claims["exp"].(float64))
			expirationDate := time.Unix(expirationTimeStamp, 0)

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
		}
	})
}
