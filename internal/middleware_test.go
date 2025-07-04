package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAppMiddleware(t *testing.T) {
	// Create a new Chi router
	r := chi.NewRouter()

	// Apply the middleware
	appMiddleware(r)

	// Define a simple test handler
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Make a request to the /ping endpoint
	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
