package api

import (
	"encoding/json"
	"net/http"
)

// OKResponse returns a JSON encoded response data with a 200 status code
func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Custom response writer with a status along side data
func Response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ErrorResponse writes a custom JSON as an error response for the client
// and if that doesn't work it will return a plain internal server error status
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	m := struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: message,
	}

	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
