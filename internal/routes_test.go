package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Setup routes with mock handlers
	r.Route("/api/v1", func(r chi.Router) {
		// Tier routes
		r.Get("/tier/all", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Fetched all tiers"))
		})
		r.Get("/tier/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Fetched tier"))
		})
		r.Post("/tier", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Tier created"))
		})
		r.Put("/tier/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Tier updated"))
		})

		// Roles routes
		r.Get("/roles/all", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("roles fetched"))
		})
		r.Get("/roles/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role fetched"))
		})
		r.Get("/roles/name", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role by name"))
		})
		r.Patch("/roles/bin", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role binned"))
		})
		r.Patch("/roles/update", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role updated"))
		})
		r.Patch("/roles/archive", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role archived"))
		})
		r.Patch("/roles/unarchive", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role unarchived"))
		})
		r.Post("/roles", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("role created"))
		})
		r.Delete("/roles/delete", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("role deleted"))
		})
	})

	// Test cases
	testCases := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		// Tier routes
		{"GET /tier/all", http.MethodGet, "/api/v1/tier/all", http.StatusOK, "Fetched all tiers"},
		{"GET /tier/{id}", http.MethodGet, "/api/v1/tier/123", http.StatusOK, "Fetched tier"},
		{"POST /tier", http.MethodPost, "/api/v1/tier", http.StatusCreated, "Tier created"},
		{"PUT /tier/{id}", http.MethodPut, "/api/v1/tier/456", http.StatusOK, "Tier updated"},

		// Roles routes
		{"GET /roles/all", http.MethodGet, "/api/v1/roles/all", http.StatusOK, "roles fetched"},
		{"GET /roles/{id}", http.MethodGet, "/api/v1/roles/789", http.StatusOK, "role fetched"},
		{"GET /roles/name", http.MethodGet, "/api/v1/roles/name", http.StatusOK, "role by name"},
		{"PATCH /roles/bin", http.MethodPatch, "/api/v1/roles/bin", http.StatusOK, "role binned"},
		{"PATCH /roles/update", http.MethodPatch, "/api/v1/roles/update", http.StatusOK, "role updated"},
		{"PATCH /roles/archive", http.MethodPatch, "/api/v1/roles/archive", http.StatusOK, "role archived"},
		{"PATCH /roles/unarchive", http.MethodPatch, "/api/v1/roles/unarchive", http.StatusOK, "role unarchived"},
		{"POST /roles", http.MethodPost, "/api/v1/roles", http.StatusCreated, "role created"},
		{"DELETE /roles/delete", http.MethodDelete, "/api/v1/roles/delete", http.StatusOK, "role deleted"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}
