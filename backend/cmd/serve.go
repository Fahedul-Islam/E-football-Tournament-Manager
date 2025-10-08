package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"tournament-manager/config"
	"tournament-manager/infra/db"
	participantrepo "tournament-manager/repository/participant_repo"
	tournamentmanagerrepo "tournament-manager/repository/tournament_manager_repo"
	userrepo "tournament-manager/repository/user-repo"
	"tournament-manager/rest/handler/participant"
	tournamentmanager "tournament-manager/rest/handler/tournamentManager"
	"tournament-manager/rest/handler/user"
	"tournament-manager/rest/middleware"
)

func Serve() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	connStr := cfg.GetDBConStr()
	var dB *sql.DB
	dB, err = db.DbConnections(connStr)
	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}
	err = db.Migrate(cfg.GetDBURL())
	if err != nil {
		fmt.Println("Database migration error:", err)
		return
	}
	mux := http.NewServeMux()

	middlewareManager := middleware.NewMiddlewareManager()
	middlewareManager.Use(middleware.Logger, middleware.CorsWithPreflight)

	userRepo := userrepo.NewUserRepo(dB)
	userHandler := user.NewUserHandler(cfg, userRepo)
	userHandler.RegisterRoutes(mux, middlewareManager)

	tournamentRepo := tournamentmanagerrepo.NewTournamentManagerRepo(dB)
	tournamentHandler := tournamentmanager.NewTournamentManagerHandler(tournamentRepo)
	tournamentHandler.RegisterRoutes(mux, middlewareManager)

	participantrepo := participantrepo.NewParticipantRepo(dB)
	participantHandler := participant.NewParticipantHandler(participantrepo)
	participantHandler.RegisterRoutes(mux, middlewareManager)

	wrappedMux := middlewareManager.WrappedMux(mux)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		fmt.Println("Server error:", err)
	}
}
