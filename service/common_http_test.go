package service

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

//jwtmiddleware "github.com/auth0/go-jwt-middleware"
//jwt "github.com/dgrijalva/jwt-go"

// path: is used to configure router path(eg: /product/{id})
// requestURL: current request path (eg: /product/1)
func makeHTTPCall(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create a http request using the given parameters
	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))

	// test recorder created for capturing apiresponses
	recorder = httptest.NewRecorder()

	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}
