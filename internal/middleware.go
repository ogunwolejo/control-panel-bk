package internal

import (
	"context"
	"control-panel-bk/util"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"strings"
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
		MaxAge:           60 * 24 * 60 * 60,
	}))
	m.Use(middleware.NoCache) // No caching
	m.Use(middleware.AllowContentEncoding("application/json", "text/xml"))
	m.Use(middleware.Compress(5, "application/json", "application/text"))
	m.Use(middleware.Heartbeat("/ping"))
	m.Use(middleware.RequestID)
	m.Use(middleware.CleanPath)
	m.Use(AppAuthorizationMiddleware)
}

func getBearerToken(r *http.Request) (*string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))

	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid Authorization header format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return nil, errors.New("empty token")
	}

	return &token, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getBearerToken(r)

		if err != nil {
			util.ErrorException(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "access_token", token)
		newReq := r.WithContext(ctx)
		next.ServeHTTP(w, newReq)
	}
}

func AppAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		request.URL.OmitHost = true
		requestMethod := request.Method
		path := request.URL.Path
		frg := request.URL.Fragment
		rr := request.URL.RequestURI()

		log.Printf("App-Auth: Request Method is: %s, and the path is: %s, Fragments are: %s, RequestURI is: %s", requestMethod, path, frg, rr)
		next.ServeHTTP(writer, request)
	})
}
