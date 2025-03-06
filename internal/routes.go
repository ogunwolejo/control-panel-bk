package internal

import (
	"control-panel-bk/pkg/tiers"
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

	mux.Route("/api/v1", func(r chi.Router) {

		r.Route("/tier", func(rtr chi.Router) {
			rtr.Post("/", tiers.HandleTierCreation)
		})

	})

	return mux
}
