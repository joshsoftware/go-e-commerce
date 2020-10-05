package service

import (
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

//PingResponse Struct
type PingResponse struct {
	Message string `json:"message"`
}

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	response := PingResponse{Message: "pong"}

	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error marshalling ping response")
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(respBytes)
}
