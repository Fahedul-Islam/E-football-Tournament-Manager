package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := dB.PingContext(r.Context()); err != nil {
			http.Error(w, `{"status":"error","db":"unreachable"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"status":"ok","db":"connected"}`)
	})

	// Initialize repositories
	userRepo := userrepo.NewUserRepo(dB)
	tournamentRepo := tournamentrepo.NewTournamentManagerRepo(dB)
	participantRepo := participantrepo.NewParticipantRepo(dB)
	announcementRepo := announcementrepo.NewAnnouncementRepo(dB)

	// Initialize services
	userSvc := userservice.NewUserService(cfg, userRepo)
	tournamentSvc := tournamentservice.NewTournamentService(tournamentRepo)
	participantSvc := participantservice.NewParticipantService(participantRepo)
	announcementSvc := announcementservice.NewAnnouncementService(announcementRepo, hub)

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

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      wrappedMux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in background goroutine
	go func() {
		fmt.Printf("Server starting on %s\n", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}
	fmt.Println("Server stopped")
}
