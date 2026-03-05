package announcement

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *AnnouncementHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
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
	announcement, err := h.announcementService.CreateAnnouncement(r.Context(), tournamentID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusCreated)
}

func (h *AnnouncementHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	announcements, err := h.announcementService.GetAnnouncements(r.Context(), tournamentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcements, http.StatusOK)
}

func (h *AnnouncementHandler) GetAnnouncementByID(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
	if err != nil {
		http.Error(w, "Invalid announcement_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	announcement, err := h.announcementService.GetAnnouncementByID(r.Context(), tournamentID, announcementID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusOK)
}

func (h *AnnouncementHandler) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
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
	announcement, err := h.announcementService.UpdateAnnouncement(r.Context(), tournamentID, announcementID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusOK)
}

func (h *AnnouncementHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
	if err != nil {
		http.Error(w, "Invalid announcement_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	err = h.announcementService.DeleteAnnouncement(r.Context(), tournamentID, announcementID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Announcement deleted successfully"}, http.StatusOK)
}

func (h *AnnouncementHandler) GetParticipantsAnnouncementSeenStatus(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
	if err != nil {
		http.Error(w, "Invalid announcement_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	status, err := h.announcementService.GetParticipantsAnnouncementSeenStatus(r.Context(), tournamentID, announcementID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, status, http.StatusOK)
}

func (h *AnnouncementHandler) ReactOnAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
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
	announcement, err := h.announcementService.ReactOnAnnouncement(r.Context(), tournamentID, announcementID, userID, reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, announcement, http.StatusOK)
}
