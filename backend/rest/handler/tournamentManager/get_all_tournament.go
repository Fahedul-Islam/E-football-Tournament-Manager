package tournamentmanager

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GetAllTournaments(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id := r.Context().Value("user_id").(string)
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	tournaments, err := h.tournamentService.GetAllTournaments(tournament_owner_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, tournaments, http.StatusOK)
}
