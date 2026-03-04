package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournament_id := r.URL.Query().Get("tournament_id")
	tournamentID, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	var req domain.AnnouncementCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	announcement, err := h.tournamentService.CreateAnnouncement(r.Context(), tournamentID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusCreated)

}

func (h *TournamentManagerHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	tournament_id := r.URL.Query().Get("tournament_id")
	tournamentID, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// only tournament owner and participants can see the announcements
	announcements, err := h.tournamentService.GetAnnouncements(r.Context(), tournamentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcements, http.StatusOK)
}

func (h *TournamentManagerHandler) GetAnnouncementByID(w http.ResponseWriter, r *http.Request) {
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

	// only tournament owner and participants can see the announcement
	announcement, err := h.tournamentService.GetAnnouncementByID(r.Context(), tournamentID, announcementID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusOK)
}

func (h *TournamentManagerHandler) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
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
	var req domain.AnnouncementCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	announcement, err := h.tournamentService.UpdateAnnouncement(r.Context(), tournamentID, announcementID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusOK)

}

func (h *TournamentManagerHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
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
	err = h.tournamentService.DeleteAnnouncement(r.Context(), tournamentID, announcementID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Announcement deleted successfully"}, http.StatusOK)
}
