package apperrors

import (
	"encoding/json"
	"errors"
	l "github.com/sirupsen/logrus"
	"net/http"
)

// ErrorStruct - struct used to convert error messages into required JSON format
type ErrorStruct struct {
	Message string `json:"message,omitempty"` //Error Message
	Status  int    `json:"status,omitempty"`  //HTTP Response status code
}

// Error - prints out an error
func Error(appError error, msg string, triggeringError error) {
	l.WithFields(l.Fields{"appError": appError, "message": msg}).Error(triggeringError)
}

// Warn - for warnings
func Warn(appError error, msg string, triggeringError error) {
	l.WithFields(l.Fields{"appError": appError, "message": msg}).Warn(triggeringError)
}

// JSONError - This function writes out an error response with the status
// header passed in
func JSONError(rw http.ResponseWriter, status int, err error) {

	errObj := ErrorStruct{
		Message: err.Error(),
		Status:  status,
	}

	errJSON, err := json.Marshal(&errObj)
	if err != nil {
		Warn(err, "Error in AppErrors marshalling JSON", err)
	}
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(errJSON)
	return
}

// ErrRecordNotFound - for when a database record isn't found
var ErrRecordNotFound = errors.New("Database record not found")

// ErrInvalidToken - used when a JSON Web Token ("JWT") cannot be validated
// by the JWT library
var ErrInvalidToken = errors.New("Invalid Token")

// ErrSignedString - failed to sign the token string
var ErrSignedString = errors.New("Failed to sign token string")

// ErrMissingAuthHeader - When the HTTP request doesn't contain an 'Authorization' header
var ErrMissingAuthHeader = errors.New("Missing Auth header")

// ErrJSONParseFail - If json.Unmarshal or json.Marshal returns an error
var ErrJSONParseFail = errors.New("Failed to parse JSON response (likely not valid JSON)")

// ErrNoSigningKey - there isn't a signing key defined in the app configuration
var ErrNoSigningKey = errors.New("no JWT signing key specified; cannot authenticate users. Define JWT_SECRET in application.yml and restart")

// ErrFailedToCreate - Record Creation Failed
var ErrFailedToCreate = errors.New("Failed to create database record")

// ErrUnknown - Generic Error For Unknown Errors
var ErrUnknown = errors.New("unknown/unexpected error has occurred")
