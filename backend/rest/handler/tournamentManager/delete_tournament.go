package tournamentmanager

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	tournament_id := r.URL.Query().Get("id")
	id, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.DeleteTournament(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Tournament deleted successfully"}, http.StatusOK)
}
