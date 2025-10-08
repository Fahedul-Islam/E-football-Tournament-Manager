package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id := r.Context().Value("user_id").(string)
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	var req domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.RemoveParticipant(tournament_owner_id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant removed successfully"}, http.StatusOK)
}
