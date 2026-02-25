package tournamentmanager

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GetAllParticipant(w http.ResponseWriter, r *http.Request) {
	str_id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(str_id)
	if err!=nil{
		http.Error(w, "Invalid tournament id", http.StatusBadRequest )
	}
	participant, err := h.tournamentService.GetAllParticipant(tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	utils.SendData(w, participant, http.StatusOK)
}
