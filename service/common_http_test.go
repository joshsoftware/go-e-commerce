package service

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"joshsoftware/go-e-commerce/config"
	"testing"

	"github.com/stretchr/testify/suite"
)

)

//jwtmiddleware "github.com/auth0/go-jwt-middleware"
//jwt "github.com/dgrijalva/jwt-go"

// path: is used to configure router path(eg: /product/{id})
// requestURL: current request path (eg: /product/1)

func TestExampleTestSuite(t *testing.T) {
	config.Load()
	suite.Run(t, new(CartHandlerTestSuite))
}

func makeHTTPCall(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	JWTToken, _ := generateJwt(1, false)

	// create a http request using the given parameters
	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))
	req.Header.Set("Token", JWTToken)

	// test recorder created for capturing api responses
	recorder = httptest.NewRecorder()
	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}

func makeHTTPCallWithJWTMiddleware(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create jwt token with userID
	JWTToken, _ := generateJwt(1, false)

	// create a http request using the given parameters
	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))
	req.Header.Set("Token", JWTToken)

	// test recorder created for capturing api responses
	recorder = httptest.NewRecorder()

	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey(), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	router.Handle(path, negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handlerFunc)),
	)).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}
func makeHTTPCallWithHeader(method, path, requestURL string, writer *multipart.Writer, payload *bytes.Buffer, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create a http request using the given parameters

	// Don't forget to set the content type, this will contain the boundary.

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	req, _ := http.NewRequest(method, requestURL, strings.NewReader(payload.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//fmt.Println(req)

	// test recorder created for capturing apiresponses
	recorder = httptest.NewRecorder()

	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}
