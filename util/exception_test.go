package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorException ensures the function correctly writes error responses.
func TestErrorException(t *testing.T) {
	// Create a mock response writer
	recorder := httptest.NewRecorder()

	// Define the test error and status code
	testErr := errors.New("something went wrong")
	errorCode := http.StatusInternalServerError

	// Call the function
	ErrorException(recorder, testErr, errorCode)

	// Verify status code
	assert.Equal(t, errorCode, recorder.Code, "Expected status code should match")

	// Verify Content-Type
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Content-Type should be application/json")

	// Verify JSON response body
	var responseBody map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "Response should be valid JSON")

	// Verify the error message
	assert.Equal(t, testErr.Error(), responseBody["error"], "Error message should match")
}
