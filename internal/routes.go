package internal

import (
	"control-panel-bk/pkg/panelAdmins"
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

			// The Panel-Admins Sub Routes
			// Role sub-router
			r.Route("/roles", func(roleRouter chi.Router) {
				roleRouter.Post("/", panelAdmins.HandleCreateRole)
				roleRouter.Get("/all", panelAdmins.HandleFetchRoles)
				roleRouter.Get("/{id}", panelAdmins.HandleFetchRoleById)
				roleRouter.Get("/name", panelAdmins.HandleFetchRoleByName) // takes the query params page and limit

				roleRouter.Put("/update", panelAdmins.HandleGeneralUpdate)
				roleRouter.Patch("/archive", panelAdmins.HandleArchiveRole)
				roleRouter.Patch("/unarchive", panelAdmins.HandleUnArchiveRole)
				roleRouter.Patch("/bin", panelAdmins.HandlePushRoleToBin)

				roleRouter.Delete("/delete", panelAdmins.HandleHardDeleteOfRole)
			})

			// Team sub-router
			r.Route("/team", func(teamRouter chi.Router) {

			})

			// User sub-router
			r.Route("/user", func(userRouter chi.Router) {

			})

		})
	})

	return mux
}
