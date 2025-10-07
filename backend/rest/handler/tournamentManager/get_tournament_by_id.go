package tournamentmanager

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GetTournamentByID(w http.ResponseWriter, r *http.Request) {
	tournament_id := r.URL.Query().Get("id")
	id, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	tournament, err := h.tournamentService.GetTournamentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tournament == nil {
		http.Error(w, "Tournament not found", http.StatusNotFound)
		return
	}
	utils.SendData(w, tournament, http.StatusOK)
}
