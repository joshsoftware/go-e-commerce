package service

import (
	ae "joshsoftware/go-e-commerce/apperrors"
	"joshsoftware/go-e-commerce/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

func TestExampleTestSuite(t *testing.T) {
	config.Load()
	suite.Run(t, new(UsersHandlerTestSuite))
	suite.Run(t, new(SessionsHandlerTestSuite))
}

// path: is used to configure router path(eg: /product/{id})
// requestURL: current request path (eg: /product/1)
func makeHTTPCall(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create a http request using the given parameters
	req, err := http.NewRequest(method, requestURL, strings.NewReader(body))
	if err != nil {
		ae.Error(ae.ErrFailedToCreate, "Error creating New request ", err)
	}
	req.Header.Set("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDU2ODQ4NjksImlkIjoxLCJpc0FkbWluIjpmYWxzZX0.-rgCTepUicCXyNLL1KiRudxT6NfowuzO2iC4oLn4reg")
	// test recorder created for capturing apiresponses
	recorder = httptest.NewRecorder()

	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}
