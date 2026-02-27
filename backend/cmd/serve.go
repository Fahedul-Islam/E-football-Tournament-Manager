package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"tournament-manager/config"
	"tournament-manager/infra/db"
	"tournament-manager/internal/delivery/http/handler/participant"
	tournamentmanager "tournament-manager/internal/delivery/http/handler/tournamentManager"
	"tournament-manager/internal/delivery/http/handler/user"
	"tournament-manager/internal/delivery/http/middleware"
	participantrepo "tournament-manager/internal/repository/participant_repo"
	tournamentmanagerrepo "tournament-manager/internal/repository/tournament_manager_repo"
	userrepo "tournament-manager/internal/repository/user-repo"

	// Service layer imports
	participantservice "tournament-manager/internal/service/participant"
	tournamentservice "tournament-manager/internal/service/tournament"
	userservice "tournament-manager/internal/service/user"
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

	// Initialize repositories
	userRepo := userrepo.NewUserRepo(dB)
	tournamentRepo := tournamentmanagerrepo.NewTournamentManagerRepo(dB)
	participantRepo := participantrepo.NewParticipantRepo(dB)

	// Initialize services
	userSvc := userservice.NewUserService(cfg, userRepo)
	tournamentSvc := tournamentservice.NewTournamentService(tournamentRepo)
	participantSvc := participantservice.NewParticipantService(participantRepo)

	// Initialize handlers with services
	userHandler := user.NewUserHandler(userSvc)
	userHandler.RegisterRoutes(mux, middlewareManager)

	tournamentHandler := tournamentmanager.NewTournamentManagerHandler(tournamentSvc)
	tournamentHandler.RegisterRoutes(mux, middlewareManager)

	participantHandler := participant.NewParticipantHandler(participantSvc)
	participantHandler.RegisterRoutes(mux, middlewareManager)

	wrappedMux := middlewareManager.WrappedMux(mux)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		fmt.Println("Server error:", err)
	}
}
