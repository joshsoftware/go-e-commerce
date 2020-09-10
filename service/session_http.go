package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ae "joshsoftware/go-e-commerce/apperrors"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
)

type OAuthToken struct {
	AccessToken string `json:"access_token"`
}

type OAuthUser struct {
	Email string `json:"email"`
	Name  string `json: "name"`
}

type authBody struct {
	Message string `json:"message"`
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
			ae.Error(ae.ErrMissingAuthHeader, "Authentication Token in Header Missing", err)
			ae.JSONError(rw, http.StatusUnauthorized, err)
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

func handleAuth(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		oauthToken := OAuthToken{}

		// Getting google access token from body
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in reading request body")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Unmarshalling access token in oauthToken
		err = json.Unmarshal(reqBody, &oauthToken)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while Unmarshalling request json")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Now getting user profile from google api using access token
		client := &http.Client{}
		req, err = http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Fail to create oauth request")
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+oauthToken.AccessToken)
		resp, err := client.Do(req)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Failure executing HTTP request to https://www.googleapis.com/oauth2/v2/userinfo", err)
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}

		u := OAuthUser{}
		payload, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error reading response body: "+string(payload), err)
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}

		err = json.Unmarshal(payload, &u)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Failure parsing JSON in Unmarshalling OAuthUser"+string(payload), err)
			ae.JSONError(rw, http.StatusInternalServerError, err)
		}

		// Checking if user is already registered if registered then generate JWT token for him
		check, existingUser, err := deps.Store.CheckUserByEmail(req.Context(), u.Email)
		if check {
			token, err := generateJwt(existingUser.ID)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Unknown/unexpected error while creating JWT for " + existingUser.Email)
				ae.JSONError(rw, http.StatusInternalServerError, err)
				return
			}
			respBody, err := json.Marshal(authBody{Message: "Authentication Successful", Token: token})
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error parsing JSON for token response, token: " + token)
				ae.JSONError(rw, http.StatusInternalServerError, err)
				return
			}
			rw.Header().Add("Content-type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(respBody)
			return
		}

		// At this point it is known that user is not registered. Now register the user and generate JWT token for him
		user := db.User{}
		user.Email = u.Email
		user.FirstName = u.Name
		newUser, err := deps.Store.CreateNewUser(req.Context(), user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error in inserting user in database")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Generating Jwt token
		token, err := generateJwt(newUser.ID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Unknown/unexpected error while creating JWT for " + newUser.Email)
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}
		respBody, err := json.Marshal(authBody{Message: "Authentication Successful", Token: token})
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error parsing JSON for token response, token: " + token)
			ae.JSONError(rw, http.StatusInternalServerError, err)
			return
		}
		rw.Header().Add("Content-type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(respBody)
		return
	})
}
