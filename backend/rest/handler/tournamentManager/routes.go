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
	mux.Handle("POST /tournaments/approve", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.ApproveParticipant)))
	mux.Handle("POST /tournaments/reject", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.RejectParticipant)))
	mux.Handle("POST /tournaments/addparticipant", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.AddParticipant)))
	mux.Handle("POST /tournaments/removeparticipant", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.RemoveParticipant)))
	mux.Handle("GET /tournaments/participants", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetAllParticipant)))
}
