package service

import (
	"joshsoftware/go-e-commerce/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// jwtmiddleware "github.com/auth0/go-jwt-middleware"
	// jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

func TestExampleTestSuite(t *testing.T) {
	config.Load()
	suite.Run(t, new(CartHandlerTestSuite))
}

func makeHTTPCall(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create a http request using the given parameters
	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))
	// test recorder created for capturing api responses
	recorder = httptest.NewRecorder()
	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}
