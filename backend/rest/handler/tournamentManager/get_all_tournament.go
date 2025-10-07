package tournamentmanager

import (
	"net/http"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GetAllTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.tournamentService.GetAllTournaments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, tournaments, http.StatusOK)
}
