package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	var participant domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.AddParticipant(participant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant added successfully"}, http.StatusOK)
}
