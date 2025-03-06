package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func appMiddleware(m *chi.Mux) {
	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Wild card
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodTrace, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodConnect, http.MethodPatch},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "X-CSRF-Token", "Host", "Origin", "Authorization", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	m.Use(middleware.AllowContentEncoding("application/json", "text/xml"))
	m.Use(middleware.Heartbeat("/ping"))
	m.Use(middleware.RequestID)
	m.Use(middleware.CleanPath)
}
