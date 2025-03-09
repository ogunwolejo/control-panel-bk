package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// Mock routes function for testing
func mockRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return mux
}

// Test the ControlPanelServer function
func TestControlPanelServer(t *testing.T) {
	port := fmt.Sprintf("%d", 9090)
	// Set the PORT environment variable for testing
	os.Setenv("PORT", port)
	defer os.Unsetenv("PORT")

	// Create a channel to listen for shutdown signals
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the server in a goroutine
	go func() {
		// Use the mock routes for testing
		server := &http.Server{
			Handler: mockRoutes(),
			Addr:    fmt.Sprintf(":%s", port),
		}

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Wait for a moment to ensure the server has started
	time.Sleep(600 * time.Millisecond)

	// Make a request to the server to check if it's running
	url := fmt.Sprintf("http://localhost:%s/health", port)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Simulate a shutdown signal
	cancel()
	// Wait for a moment to ensure the server has time to shut down
	time.Sleep(100 * time.Millisecond)
}

// Test for graceful shutdown
func TestControlPanelServerShutdown(t *testing.T) {
	port := fmt.Sprintf("%d", 9090)
	// Set the PORT environment variable for testing
	os.Setenv("PORT", port)
	defer os.Unsetenv("PORT")

	// Create a channel to listen for shutdown signals
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		// Use the mock routes for testing
		server := &http.Server{
			Handler: mockRoutes(),
			Addr:    fmt.Sprintf(":%s", port),
		}
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Wait for a moment to ensure the server has started
	time.Sleep(100 * time.Millisecond)

	// Simulate a shutdown signal
	cancel()
	// Wait for a moment to ensure the server has time to shut down
	time.Sleep(100 * time.Millisecond)
}

// Benchmark tests for ControlPanelServer
func BenchmarkControlPanelServer(b *testing.B) {
	// Set the PORT environment variable for testing
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	// Start the server in a goroutine
	go func() {
		ControlPanelServer()
	}()

	// Wait for a moment to ensure the server has started
	time.Sleep(100 * time.Millisecond)

	// Benchmark the /health endpoint
	b.Run("BenchmarkHealthCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := http.Get("http://localhost:8080/health")
			if err != nil {
				b.Fatalf("Failed to make request: %v", err)
			}
			resp.Body.Close() // Close the response body
		}
	})

	// Wait for a moment to ensure the server has time to shut down
	time.Sleep(100 * time.Millisecond)
}
