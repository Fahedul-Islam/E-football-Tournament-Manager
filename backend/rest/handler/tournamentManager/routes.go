package tournamentmanager

import (
	"net/http"
	"tournament-manager/rest/middleware"
)

func (h *TournamentManagerHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("POST /tournaments/create", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.CreateTournament)))
	mux.Handle("GET /tournaments", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetAllTournaments)))
	mux.Handle("DELETE /tournaments", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.DeleteTournament)))
	mux.Handle("PUT /tournaments", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.UpdateTournament)))
	mux.Handle("PATCH /tournaments/approve", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.ApproveParticipant)))
	mux.Handle("PATCH /tournaments/reject", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.RejectParticipant)))
	mux.Handle("POST /tournaments/addparticipant", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.AddParticipant)))
	mux.Handle("POST /tournaments/removeparticipant", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.RemoveParticipant)))
	mux.Handle("GET /tournaments/participants", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetAllParticipant)))
	mux.Handle("GET /tournaments/create_match_schedules", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.CreateMatchSchedules)))
	mux.Handle("GET /tournaments/matches", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetAllMatches)))
	mux.Handle("PATCH /tournaments/matche-score/update",manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.UpdateScore)))
	
}
