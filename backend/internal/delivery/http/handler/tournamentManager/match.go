package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

// CreateMatchSchedules creates match schedules for a group+knockout type tournament
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
	// handling other tournament types
	cnt := r.URL.Query().Get("group_count")
	groupCount, err := strconv.Atoi(cnt)

	if err != nil {
		http.Error(w, "Invalid group count", http.StatusBadRequest)
		return
	}
	err = h.tournamentService.CreateMatchSchedules(r.Context(), tournament_id, tournament_owner_id, groupCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// LeagueStyleSchedule creates a league style schedule for a tournament
func (h *TournamentManagerHandler) LeagueStyleSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}

	err = h.tournamentService.LeagueStyleSchedule(r.Context(), tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Update Score updates the score of a match and advances the tournament if necessary
func (h *TournamentManagerHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id := r.Context().Value("user_id").(string)
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var req domain.UpdateMatchScoreInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.tournamentService.UpdateScore(r.Context(), tournament_owner_id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, result, http.StatusOK)
}

// Get all the matches of a tournament
func (h *TournamentManagerHandler) GetAllMatches(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	matches, err := h.tournamentService.GetAllMatches(r.Context(), tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, matches, http.StatusOK)
}
