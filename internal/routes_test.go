package internal

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type RoutesTestSuite struct {
	suite.Suite
	router *chi.Mux
}

func mockHandler(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}
}

// Run the test suite
func TestRoutesSuite(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}

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

func (suite *RoutesTestSuite) SetupSuite() {
	suite.router = chi.NewRouter()

	// Define mock team routes
	suite.router.Route("/teams", func(teamRouter chi.Router) {
		teamRouter.Post("/create", mockHandler(http.StatusCreated, "Team created"))

		teamRouter.Patch("/archive", mockHandler(http.StatusOK, "Team archived"))
		teamRouter.Patch("/unarchive", mockHandler(http.StatusOK, "Team unarchived"))
		teamRouter.Patch("/add-members", mockHandler(http.StatusOK, "Members added"))
		teamRouter.Patch("/remove-members", mockHandler(http.StatusOK, "Members removed"))
		teamRouter.Patch("/change-lead", mockHandler(http.StatusOK, "Team lead changed"))
		teamRouter.Patch("/bin", mockHandler(http.StatusOK, "Team moved to bin"))
		teamRouter.Patch("/restore", mockHandler(http.StatusOK, "Team restored"))

		teamRouter.Delete("/delete", mockHandler(http.StatusOK, "Team deleted"))

		teamRouter.Get("/{id}", mockHandler(http.StatusOK, "Team details"))
		teamRouter.Get("/all", mockHandler(http.StatusOK, "All teams"))
	})


	suite.router.Route("/users", func(userRouter chi.Router) {})
	suite.router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/create", mockHandler(http.StatusCreated, "New User was created"))
		authRouter.Post("/login", mockHandler(http.StatusOK, "User is login"))
		authRouter.Post("/logout", mockHandler(http.StatusOK, "User was logout"))
		authRouter.Get("/refresh-token", mockHandler(http.StatusOK, "Token regenerated"))
	})
}

// Test cases
func (suite *RoutesTestSuite) TestTeamRoutes() {
	testCases := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		// Create team
		{"POST /teams/create", http.MethodPost, "/teams/create", http.StatusCreated, "Team created"},

		// Modify team
		{"PATCH /teams/archive", http.MethodPatch, "/teams/archive", http.StatusOK, "Team archived"},
		{"PATCH /teams/unarchive", http.MethodPatch, "/teams/unarchive", http.StatusOK, "Team unarchived"},
		{"PATCH /teams/add-members", http.MethodPatch, "/teams/add-members", http.StatusOK, "Members added"},
		{"PATCH /teams/remove-members", http.MethodPatch, "/teams/remove-members", http.StatusOK, "Members removed"},
		{"PATCH /teams/change-lead", http.MethodPatch, "/teams/change-lead", http.StatusOK, "Team lead changed"},
		{"PATCH /teams/bin", http.MethodPatch, "/teams/bin", http.StatusOK, "Team moved to bin"},
		{"PATCH /teams/restore", http.MethodPatch, "/teams/restore", http.StatusOK, "Team restored"},

		// Delete team
		{"DELETE /teams/delete", http.MethodDelete, "/teams/delete", http.StatusOK, "Team deleted"},

		// Get team details
		{"GET /teams/{id}", http.MethodGet, "/teams/123", http.StatusOK, "Team details"},
		{"GET /teams/all", http.MethodGet, "/teams/all", http.StatusOK, "All teams"},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			suite.router.ServeHTTP(rec, req)

			assert.Equal(suite.T(), tc.wantStatus, rec.Code)
			assert.Contains(suite.T(), rec.Body.String(), tc.wantBody)
		})
	}
}

func (suite *RoutesTestSuite) TestAuthRoutes() {
	testCases := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		// Create team
		{"POST /auth/create", http.MethodPost, "/auth/create", http.StatusCreated, "New User was created"},
		{"POST /auth/login", http.MethodPost, "/auth/login", http.StatusOK, "User is login"},
		{"POST /auth/logout", http.MethodPost, "/auth/logout", http.StatusOK, "User was logout"},

		// Get team details
		{"GET /auth/refresh-token", http.MethodGet, "/auth/refresh-token", http.StatusOK, "Token regenerated"},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			suite.router.ServeHTTP(rec, req)

			assert.Equal(suite.T(), tc.wantStatus, rec.Code)
			assert.Contains(suite.T(), rec.Body.String(), tc.wantBody)
		})
	}
}
