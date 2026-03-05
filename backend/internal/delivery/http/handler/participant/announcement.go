package participant

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *ParticipantHandler) ReactOnAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournament_id := r.URL.Query().Get("tournament_id")
	tournamentID, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcement_id := r.URL.Query().Get("announcement_id")
	announcementID, err := strconv.Atoi(announcement_id)
	if err != nil {
		http.Error(w, "Invalid announcement_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	reaction := r.URL.Query().Get("reaction")
	if reaction != "like" && reaction != "dislike" {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}
	announcement, err := h.service.ReactOnAnnouncement(r.Context(), tournamentID, announcementID, userID, reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusOK)
}
