package participant

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *ParticipantHandler) SeeMatchSchedule(w http.ResponseWriter, r *http.Request) {
	str_user_id := r.Context().Value("user_id").(string)
	user_id, err := strconv.Atoi(str_user_id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	isApproved, err := h.service.IsApprovedParticipant(tournament_id, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isApproved {
		http.Error(w, "User is not approved", http.StatusForbidden)
		return
	}
	
	match_schedule, err := h.service.SeeMatchSchedule(tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, match_schedule, http.StatusOK)
}
