package ws

import (
	"net/http"
	"tournament-manager/internal/delivery/http/middleware"
)

func (h *WebSocketHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("/ws", manager.With(middleware.AuthMiddleware(""))(http.HandlerFunc(h.HandleWebSocket)))
}
