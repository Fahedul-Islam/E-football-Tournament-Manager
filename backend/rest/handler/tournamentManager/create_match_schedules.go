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
	// get tournament owner id to verify permission
	str_t_owner_id := r.Context().Value("user_id").(string)
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	// verify permission
	hasPermission, err := h.tournamentService.VerifyTournamentOwner(tournament_id, tournament_owner_id)
	if err != nil {
		http.Error(w, "Failed to verify tournament owner", http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		http.Error(w, "You do not have permission to create match schedules for this tournament", http.StatusForbidden)
		return
	}

	// check tournament type
	tournment_type, err := h.tournamentService.GetTournamentType(tournament_id)
	if err != nil {
		http.Error(w, "Failed to get tournament type", http.StatusInternalServerError)
		return
	}
	// handling league style tournament
	if tournment_type == "league" {
		h.LeagueStyleSchedule(w, r)
		return
	}
	// handling other tournament types
	cnt := r.URL.Query().Get("group_count")
	groupCount, err := strconv.Atoi(cnt)

	if err != nil {
		http.Error(w, "Invalid group count", http.StatusBadRequest)
		return
	}
	if (groupCount < 1 || groupCount> 8 || groupCount%2 != 0) {
		http.Error(w, "Group count must be between 1 and 8 and an even number", http.StatusBadRequest)
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
	err = h.tournamentService.CreateMatchSchedules(tournament_id, groupCount, approvedParticipants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *TournamentManagerHandler) LeagueStyleSchedule(w http.ResponseWriter, r *http.Request) {
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
	err = h.tournamentService.LeagueStyleSchedule(tournament_id, approvedParticipants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}