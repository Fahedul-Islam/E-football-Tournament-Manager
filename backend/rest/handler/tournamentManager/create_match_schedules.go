package tournamentmanager

import (
	"net/http"
	"strconv"
	"tournament-manager/domain"
)

func (h *TournamentManagerHandler) CreateMatchSchedules(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	var approvedParticipants []*domain.Participant
	approvedParticipants, err = h.tournamentService.GetApprovedParticipants(tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(approvedParticipants) < 2 {
		http.Error(w, "Not enough approved participants to create match schedules", http.StatusBadRequest)
		return
	}
	err = h.tournamentService.CreateMatchSchedules(tournament_id, approvedParticipants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
