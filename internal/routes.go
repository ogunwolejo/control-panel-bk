package internal

import (
	"control-panel-bk/internal/aws"
	"control-panel-bk/pkg"
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
			// Auth Sub Routes
			r.Route("/auth", func(authRouter chi.Router) {
				authRouter.Post("/create", panelAdmins.CreateUser(aws.MongoDBClient))
				authRouter.Get("/refresh-token", pkg.RefreshTokenAuth)
				authRouter.Post("/login", pkg.LoginHandler)
				authRouter.Get("/logout", AuthMiddleware(pkg.LogoutHandler))
				authRouter.Post("/change-password", AuthMiddleware(pkg.ChangePasswordHandle))
				authRouter.Post("/forget-password-otp", pkg.ForgetPasswordOtpHandle)
				authRouter.Post("/forget-password", pkg.ForgetPasswordHandle)
			})

			// The Tier Sub Routes
			r.Route("/tier", func(tierRouter chi.Router) {
				tierRouter.Get("/all", AuthMiddleware(tiers.HandleFetchTiers))
				tierRouter.Get("/{id}", AuthMiddleware(tiers.HandleFetchTier))

				tierRouter.Group(func(tierRouterGroup chi.Router) {
					tierRouterGroup.Post("/", AuthMiddleware(tiers.HandleTierCreation))
					tierRouterGroup.Put("/{id}", AuthMiddleware(tiers.HandleUpdateTier))
				})
			})

			// The Panel-Admins Sub Routes
			// Role sub-router
			r.Route("/roles", func(roleRouter chi.Router) {
				roleRouter.Post("/", AuthMiddleware(panelAdmins.HandleCreateRole(db)))
				roleRouter.Get("/all", AuthMiddleware(panelAdmins.HandleFetchRoles(db)))
				roleRouter.Get("/{id}", AuthMiddleware(panelAdmins.HandleFetchRoleById(db)))
				roleRouter.Get("/name", AuthMiddleware(panelAdmins.HandleFetchRoleByName(db))) // takes the query params page and limit

				roleRouter.Patch("/update", AuthMiddleware(panelAdmins.HandleGeneralUpdate(db)))
				roleRouter.Patch("/archive", AuthMiddleware(panelAdmins.HandleArchiveRole(db)))
				roleRouter.Patch("/unarchive", AuthMiddleware(panelAdmins.HandleUnArchiveRole(db)))
				roleRouter.Patch("/bin", AuthMiddleware(panelAdmins.HandlePushRoleToBin(db)))
				roleRouter.Patch("/restore", AuthMiddleware(panelAdmins.HandleRestoreRoleFromBin(db)))

				roleRouter.Delete("/delete", panelAdmins.HandleHardDeleteOfRole(db))
			})

			// Team sub-router
			r.Route("/teams", func(teamRouter chi.Router) {
				teamRouter.Post("/create", AuthMiddleware(panelAdmins.HandleCreateTeam(db)))

				teamRouter.Patch("/archive", AuthMiddleware(panelAdmins.HandleArchiveTeam(db)))
				teamRouter.Patch("/unarchive", AuthMiddleware(panelAdmins.HandleUnArchiveTeam(db)))
				teamRouter.Patch("/add-members", AuthMiddleware(panelAdmins.HandleAddNewMembers(db)))
				teamRouter.Patch("/remove-members", AuthMiddleware(panelAdmins.HandleRemoveNewMembers(db)))
				teamRouter.Patch("/change-lead", AuthMiddleware(panelAdmins.HandleChangeTeamLead(db)))
				teamRouter.Patch("/bin", AuthMiddleware(panelAdmins.PushTeamToBin(db)))
				teamRouter.Patch("/restore", AuthMiddleware(panelAdmins.RestoreTeamFromBin(db)))

				teamRouter.Delete("/delete", AuthMiddleware(panelAdmins.HardDeleteTeam(db)))

				teamRouter.Get("/{id}", AuthMiddleware(panelAdmins.GetTeam(db)))
				teamRouter.Get("/all", AuthMiddleware(panelAdmins.GetTeams(db)))
			})

			// User sub-router
			r.Route("/users", func(userRouter chi.Router) {
				userRouter.Get("/", AuthMiddleware(panelAdmins.GetUsers(db)))
				userRouter.Get("/{user}", AuthMiddleware(panelAdmins.GetUser(db)))

				userRouter.Patch("/de-active", AuthMiddleware(panelAdmins.DeActiveUser(db)))
				userRouter.Patch("/reactive", AuthMiddleware(panelAdmins.ActiveUser(db)))
			})

		})
	})

	return mux
}
