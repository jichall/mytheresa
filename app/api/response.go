package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// RespondOK returns a JSON encoded response data with a 200 status code
func RespondOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		RespondError(w, Response{Status: http.StatusInternalServerError, Error: err.Error()})
	}
}

// Custom response writer with a status along side data
func RespondCustom(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RespondError writes a custom JSON as an error response for the client
// and if that doesn't work it will return a plain internal server error status
func RespondError(w http.ResponseWriter, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.Status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
