package participant

import (
	"net/http"
	"tournament-manager/rest/middleware"
)

func (h *ParticipantHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("POST /join-tournament", manager.With()(http.HandlerFunc(h.RequestToJoin)))
}
