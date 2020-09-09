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

//InitRouter :The routing mechanism. Mux helps us define handler functions and the access methods
func InitRouter(deps Dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())
	fmt.Println(v1)

	router.HandleFunc("/users", listUsersHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/users/{id:[0-9]+}", getUserHandler(deps)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/users/{id:[0-9]+}", updateUserHandler(deps)).Methods(http.MethodPut).Headers(versionHeader, v1)
	return
}
