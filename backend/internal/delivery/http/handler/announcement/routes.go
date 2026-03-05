package announcement

import (
	"net/http"
	"tournament-manager/internal/delivery/http/middleware"
)

func (h *AnnouncementHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	// Admin routes for announcement management
	mux.Handle("POST /tournaments/announcements", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.CreateAnnouncement)))
	mux.Handle("PUT /tournaments/announcements/update", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.UpdateAnnouncement)))
	mux.Handle("DELETE /tournaments/announcements/delete", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.DeleteAnnouncement)))
	mux.Handle("GET /tournaments/announcements/seen_status", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetParticipantsAnnouncementSeenStatus)))

	// Routes for both admin and participants
	mux.Handle("GET /tournaments/announcements", manager.With(middleware.AuthMiddleware(""))(http.HandlerFunc(h.GetAnnouncements)))
	mux.Handle("GET /tournaments/announcements/get", manager.With(middleware.AuthMiddleware(""))(http.HandlerFunc(h.GetAnnouncementByID)))

	// Participant routes for reactions
	mux.Handle("POST /tournament/announcement/react", manager.With(middleware.AuthMiddleware("player"))(http.HandlerFunc(h.ReactOnAnnouncement)))
}
