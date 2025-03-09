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
	mux.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			// The Tier Sub Routes
			r.Route("/tier", func(tierRouter chi.Router) {
				tierRouter.Get("/all", tiers.HandleFetchTiers)
				tierRouter.Get("/{id}", tiers.HandleFetchTier)

				tierRouter.Group(func(tierRouterGroup chi.Router) {
					tierRouterGroup.Post("/", tiers.HandleTierCreation)
					tierRouterGroup.Put("/{id}", tiers.HandleUpdateTier)
				})
			})

		})
	})

	return mux
}
