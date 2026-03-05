package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOKResponse(t *testing.T) {
	sample := struct {
		Message string `json:"message"`
	}{Message: "Success"}

	t.Run("sucessfull http 200 json response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		RespondOK(recorder, sample)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"message":"Success"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}

func TestResponse(t *testing.T) {
	sample := struct {
		Message string `json:"message"`
	}{Message: "Created"}

	t.Run("sucessfull http 201 json response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Response(recorder, http.StatusCreated, sample)

		assert.Equal(t, http.StatusCreated, recorder.Code, "Expected status code 201")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"message":"Created"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}

func TestErrorResponse(t *testing.T) {
	t.Run("json response for a given http status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		RespondError(recorder, http.StatusInternalServerError, "Some error occurred")

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"error":"Some error occurred"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}
