package tournamentmanager

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GetAllMatches(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	matches, err := h.tournamentService.GetAllMatches(tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, matches, http.StatusOK)
}
