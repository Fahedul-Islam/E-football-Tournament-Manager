package tournament

import (
	"net/http"
	"strconv"
	"tournament-manager/internal/delivery/http/middleware"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GenerateGroups(w http.ResponseWriter, r *http.Request) {
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

	groupCountStr := r.URL.Query().Get("group_count")
	groupCount, err := strconv.Atoi(groupCountStr)
	if err != nil {
		http.Error(w, "Invalid group count", http.StatusBadRequest)
		return
	}
	err = h.tournamentService.GenerateGroups(r.Context(), tournament_id, groupCount, tournament_owner_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Groups generated successfully"}, http.StatusOK)
}
