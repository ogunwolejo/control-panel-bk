package internal

import (
	"bytes"
	"control-panel-bk/pkg/tiers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Mock handler for testing
func mockHandleFetchTiers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fetched all tiers"))
}

func mockHandleFetchTier(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fetched tier"))
}

func mockHandleTierCreation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Tier created"))
}

func mockHandleUpdateTier(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tier updated"))
}

// Test the routes
func TestRoutes(t *testing.T) {
	// Create a new chi router
	mux := chi.NewRouter()

	// Apply middleware (if any)
	appMiddleware(mux)

	// Define the routes with mock handlers
	mux.Route("/api/v1/tier", func(r chi.Router) {
		r.Get("/all", mockHandleFetchTiers)
		r.Get("/{id}", mockHandleFetchTier)
		r.Post("/", mockHandleTierCreation)
		r.Put("/{id}", mockHandleUpdateTier)
	})

	// Create a test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Test the /api/v1/tier/all endpoint
	resp, err := http.Get(ts.URL + "/api/v1/tier/all")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Test the /api/v1/tier/{id} endpoint
	resp, err = http.Get(ts.URL + "/api/v1/tier/1")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Test the POST /api/v1/tier/ endpoint
	resp, err = http.Post(ts.URL+"/api/v1/tier", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	// Test the PUT /api/v1/tier/{id} endpoint
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/tier/1", nil)
	if err != nil {
		t.Fatalf("Failed to create PUT request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

// Benchmark tests for routes
func BenchmarkRoutes(b *testing.B) {
	// Create a new chi router
	mux := chi.NewRouter()

	// Apply middleware (if any)
	appMiddleware(mux)

	// Define the routes with mock handlers
	mux.Route("/api/v1/tier", func(r chi.Router) {
		r.Get("/all", mockHandleFetchTiers)
		r.Get("/{id}", mockHandleFetchTier)
		r.Post("/", mockHandleTierCreation)
		r.Put("/{id}", mockHandleUpdateTier)
	})

	// Create a test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Benchmark the /api/v1/tier/all endpoint
	b.Run("BenchmarkFetchTiers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := http.Get(ts.URL + "/api/v1/tier/all")
			if err != nil {
				b.Fatalf("Failed to make request: %v", err)
			}
			resp.Body.Close() // Close the response body
		}
	})

	// Benchmark the /api/v1/tier/{id} endpoint
	b.Run("BenchmarkFetchTier", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := http.Get(ts.URL + "/api/v1/tier/1")
			if err != nil {
				b.Fatalf("Failed to make request: %v", err)
			}

			resp.Body.Close() // Close the response body
		}
	})

	// Benchmark for creating tier /api/v1/tier
	b.Run("BenchmarkCreateTier", func(b *testing.B) {
		bdy := tiers.CreateTierRequest{
			Amount:   500,
			Interval: tiers.IntervalAnnually,
			Name:     "Benchmark test",
			Currency: tiers.CurrencyNGN,
		}

		requestBody, _ := json.Marshal(bdy)
		for i := 0; i < b.N; i++ {
			resp, err := http.Post(ts.URL+"/api/v1/tier", "application/json", bytes.NewReader(requestBody))
			if err != nil {
				b.Fatalf("Failed to make POST request: %v", err)
			}
			resp.Body.Close() // Close the response body
		}
	})

	// Benchmark the PUT /api/v1/tier/{id} endpoint
	b.Run("BenchmarkUpdateTier", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/tier/1", nil)
			if err != nil {
				b.Fatalf("Failed to create PUT request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				b.Fatalf("Failed to make PUT request: %v", err)
			}
			resp.Body.Close() // Close the response body
		}
	})
}
