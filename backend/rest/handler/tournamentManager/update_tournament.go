package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) UpdateTournament(w http.ResponseWriter, r *http.Request) {
	var req domain.TournamentCreateRequest
	tournament_id := r.URL.Query().Get("id")
	id, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.UpdateTournament(id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Tournament updated successfully"}, http.StatusOK)
}
