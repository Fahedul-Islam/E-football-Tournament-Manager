package participant

import (
	"net/http"
	"tournament-manager/rest/middleware"
)

func (h *ParticipantHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("POST /join-tournament", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.RequestToJoin)))
	mux.Handle("GET /tournament/group-distribution", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.GetGroupDistribution)))
	mux.Handle("GET /tournament/match-schedule", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.SeeMatchSchedule)))
}
