package user

import (
	"net/http"
	"tournament-manager/internal/delivery/http/middleware"
)

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("POST /register", manager.With(middleware.RateLimit)(http.HandlerFunc(h.RegisterUser)))
	mux.Handle("POST /login", manager.With(middleware.RateLimit)(http.HandlerFunc(h.LoginUser)))
	mux.Handle("GET /users", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetUsers)))
}
