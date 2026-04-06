package tournament

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/delivery/http/middleware"
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
	str_t_owner_id, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
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
	utils.SendData(w, map[string]string{"message": "Match schedules created successfully"}, http.StatusOK)
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
	utils.SendData(w, map[string]string{"message": "League schedule created successfully"}, http.StatusOK)
}

// UpdateScore updates the score of a match and advances the tournament if necessary
func (h *TournamentManagerHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

// GetAllMatches returns all matches of a tournament
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
