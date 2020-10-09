package service

import (
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

type errorResponse struct {
	Error interface{} `json:"error"`
}
type successResponse struct {
	Data interface{} `json:"data"`
}

type messageObject struct {
	Message string `json:"message"`
}

func responseMsg(rw http.ResponseWriter, status int, msgbody string) {
	response(rw, status, errorResponse{
		Error: messageObject{
			Message: msgbody,
		},
	})
}

func response(rw http.ResponseWriter, status int, responseBody interface{}) {
	respBytes, err := json.Marshal(responseBody)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while marshaling core values data")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}
