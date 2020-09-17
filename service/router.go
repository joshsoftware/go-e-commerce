package service

import (
	"fmt"
	"net/http"

	"joshsoftware/go-e-commerce/config"

	"github.com/gorilla/mux"
)

const (
	versionHeader = "Accept"
)

/* The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter(deps Dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	router.HandleFunc("/users", listUsersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.HandleFunc("/user/{id:[0-9]+}", deleteUserHandler(deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.HandleFunc("/user/disable/{id:[0-9]+}", disableUserHandler(deps)).Methods(http.MethodPatch).Headers(versionHeader, v1)

	return
}
