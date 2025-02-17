package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	appMiddleware(mux)

	// Routes
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	return mux
}
