package participant

import (
	"net/http"
	"tournament-manager/internal/delivery/http/middleware"
)

func (h *ParticipantHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("POST /join-tournament", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.RequestToJoin)))
	mux.Handle("GET /tournament/group-distribution", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.GetGroupDistribution)))
	mux.Handle("GET /tournament/match-schedule", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.SeeMatchSchedule)))
	mux.Handle("POST /tournament/announcement/react", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.ReactOnAnnouncement)))
}
