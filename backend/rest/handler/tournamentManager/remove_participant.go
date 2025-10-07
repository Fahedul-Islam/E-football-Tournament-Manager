package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
	var req domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.RemoveParticipant(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant removed successfully"}, http.StatusOK)
}
