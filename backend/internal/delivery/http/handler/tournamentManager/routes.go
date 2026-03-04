package tournamentmanager

import (
	"net/http"
	"tournament-manager/internal/delivery/http/middleware"
)

func (h *TournamentManagerHandler) RegisterRoutes(mux *http.ServeMux, manager *middleware.MiddlewareManager) {
	mux.Handle("POST /tournaments/create", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.CreateTournament)))
	mux.Handle("GET /tournaments", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.AllTournaments)))
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
	mux.Handle("GET /tournaments/leaderboard",manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.GetGroupStageLeaderboard)))

	//Announcement routes
	mux.Handle("POST /tournaments/announcements", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.CreateAnnouncement)))
	mux.Handle("GET /tournaments/announcements", manager.With(middleware.AuthMiddleware(""))(http.HandlerFunc(h.GetAnnouncements)))
	mux.Handle("GET /tournaments/announcements/get", manager.With(middleware.AuthMiddleware(""))(http.HandlerFunc(h.GetAnnouncementByID)))
	mux.Handle("PUT /tournaments/announcements/update", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.UpdateAnnouncement)))
	mux.Handle("DELETE /tournaments/announcements/delete", manager.With(middleware.AuthMiddleware("admin"))(http.HandlerFunc(h.DeleteAnnouncement)))
}
