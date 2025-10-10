package participant

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *ParticipantHandler) GetGroupDistribution(w http.ResponseWriter, r *http.Request) {
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

	group_teams, err := h.service.GetGroupDistribution(tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(group_teams) == 0 {
		http.Error(w, "No group distribution found", http.StatusNotFound)
		return
	}
	utils.SendData(w, group_teams, http.StatusOK)
}
