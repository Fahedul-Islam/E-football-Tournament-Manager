package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"tournament-manager/config"
	"tournament-manager/infra/db"
	infraws "tournament-manager/infra/ws"
	"tournament-manager/internal/delivery/http/handler/announcement"
	"tournament-manager/internal/delivery/http/handler/participant"
	"tournament-manager/internal/delivery/http/handler/tournament"
	"tournament-manager/internal/delivery/http/handler/user"
	ws "tournament-manager/internal/delivery/http/handler/ws"
	"tournament-manager/internal/delivery/http/middleware"
	announcementrepo "tournament-manager/internal/repository/announcement"
	participantrepo "tournament-manager/internal/repository/participant_repo"
	tournamentrepo "tournament-manager/internal/repository/tournament_manager_repo"
	userrepo "tournament-manager/internal/repository/user"

	// Service layer imports
	announcementservice "tournament-manager/internal/service/announcement"
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

	hub := infraws.NewHub()
	go hub.Run()

	mux := http.NewServeMux()

	middlewareManager := middleware.NewMiddlewareManager()
	middlewareManager.Use(middleware.Logger, middleware.CorsWithPreflight)

	// Initialize repositories
	userRepo := userrepo.NewUserRepo(dB)
	tournamentRepo := tournamentrepo.NewTournamentManagerRepo(dB)
	participantRepo := participantrepo.NewParticipantRepo(dB)
	announcementRepo := announcementrepo.NewAnnouncementRepo(dB)

	// Initialize services
	userSvc := userservice.NewUserService(cfg, userRepo)
	tournamentSvc := tournamentservice.NewTournamentService(tournamentRepo)
	participantSvc := participantservice.NewParticipantService(participantRepo)
	announcementSvc := announcementservice.NewAnnouncementService(announcementRepo,hub)

	// Initialize handlers with services
	userHandler := user.NewUserHandler(userSvc)
	userHandler.RegisterRoutes(mux, middlewareManager)

	tournamentHandler := tournament.NewTournamentManagerHandler(tournamentSvc)
	tournamentHandler.RegisterRoutes(mux, middlewareManager)

	participantHandler := participant.NewParticipantHandler(participantSvc)
	participantHandler.RegisterRoutes(mux, middlewareManager)

	announcementHandler := announcement.NewAnnouncementHandler(announcementSvc)
	announcementHandler.RegisterRoutes(mux, middlewareManager)

	// WebSocket handler
	wsHandler := ws.WebSocketHandler{Hub: hub}
	wsHandler.RegisterRoutes(mux, middlewareManager)

	wrappedMux := middlewareManager.WrappedMux(mux)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		fmt.Println("Server error:", err)
	}
}
