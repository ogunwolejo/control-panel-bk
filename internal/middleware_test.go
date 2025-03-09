package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAppMiddleware(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Apply the middleware
	appMiddleware(r)

	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Test the /ping endpoint
	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Test a CORS preflight request
	req, err := http.NewRequest(http.MethodOptions, ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Origin", "http://example.com")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make OPTIONS request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for OPTIONS request, got %v", resp.Status)
	}
}



func BenchmarkAppMiddleware(b *testing.B) {
	// Create a new chi router
	r := chi.NewRouter()

	// Apply the middleware
	appMiddleware(r)

	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Benchmark the /ping endpoint
	b.ResetTimer() // Reset the timer before the benchmark loop
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/ping")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close() // Close the response body
	}
}