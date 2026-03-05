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

	// to avoid writing newline characters at the end of the response if data is nil,
	// as this method is called we might spend a lot of network traffic to just send
	// useless data.
	//
	// the best way though would be to implement an encoder that doesn't write a newline
	// character at the end of the response, but this is beyond what's being asked.
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
