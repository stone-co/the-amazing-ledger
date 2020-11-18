package responses

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func Send(w http.ResponseWriter, log *logrus.Entry, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.WithError(err).Info("can't write response")
	}
}

func SendError(w http.ResponseWriter, log *logrus.Entry, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{Message: errorMessage}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.WithError(err).Info("can't write response")
	}
}
