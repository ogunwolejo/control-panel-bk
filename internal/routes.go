package internal

import (
	"control-panel-bk/internal/aws"
	"control-panel-bk/pkg/panelAdmins"
	"control-panel-bk/pkg/tiers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func getDB(client *mongo.Client) *mongo.Database {
	if client == nil {
		panic("MongoDB database not initialized")
	}
	return client.Database("flowCx")
}

func Routes() *chi.Mux {
	mux := chi.NewRouter()
	appMiddleware(mux)

	db := getDB(aws.MongoDBClient)

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
				roleRouter.Post("/", panelAdmins.HandleCreateRole(db))
				roleRouter.Get("/all", panelAdmins.HandleFetchRoles(db))
				roleRouter.Get("/{id}", panelAdmins.HandleFetchRoleById(db))
				roleRouter.Get("/name", panelAdmins.HandleFetchRoleByName(db)) // takes the query params page and limit

				roleRouter.Patch("/update", panelAdmins.HandleGeneralUpdate(db))
				roleRouter.Patch("/archive", panelAdmins.HandleArchiveRole(db))
				roleRouter.Patch("/unarchive", panelAdmins.HandleUnArchiveRole(db))
				roleRouter.Patch("/bin", panelAdmins.HandlePushRoleToBin(db))
				roleRouter.Patch("/restore", panelAdmins.HandleRestoreRoleFromBin(db))

				roleRouter.Delete("/delete", panelAdmins.HandleHardDeleteOfRole(db))
			})

			// Team sub-router
			r.Route("/teams", func(teamRouter chi.Router) {
				teamRouter.Post("/create", panelAdmins.HandleCreateTeam(db))

				teamRouter.Patch("/archive", panelAdmins.HandleArchiveTeam(db))
				teamRouter.Patch("/unarchive", panelAdmins.HandleUnArchiveTeam(db))
				teamRouter.Patch("/add-members", panelAdmins.HandleAddNewMembers(db))
				teamRouter.Patch("/remove-members", panelAdmins.HandleRemoveNewMembers(db))
				teamRouter.Patch("/change-lead", panelAdmins.HandleChangeTeamLead(db))
				teamRouter.Patch("/bin", panelAdmins.PushTeamToBin(db))
				teamRouter.Patch("/restore", panelAdmins.RestoreTeamFromBin(db))

				teamRouter.Delete("/delete", panelAdmins.HardDeleteTeam(db))

				teamRouter.Get("/{id}", panelAdmins.GetTeam(db))
				teamRouter.Get("/all", panelAdmins.GetTeams(db))
			})

			// User sub-router
			r.Route("/users", func(userRouter chi.Router) {
				userRouter.Post("/create", panelAdmins.CreateUser(aws.MongoDBClient))
			})

		})
	})

	return mux
}
