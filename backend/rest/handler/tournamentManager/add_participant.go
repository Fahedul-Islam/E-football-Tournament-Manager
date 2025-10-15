package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(string)
	tournament_owner_id, err:= strconv.Atoi(id)
	if err!=nil{
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	var participant domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.AddParticipant(tournament_owner_id,participant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant added successfully"}, http.StatusOK)
}
