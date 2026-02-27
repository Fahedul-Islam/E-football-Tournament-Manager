package participant

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *ParticipantHandler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	var data domain.ParticipantRequest
	var req struct {
		TournamentID int    `json:"tournament_id"`
		TeamName     string `json:"team_name"`
	}
	str_user_id := r.Context().Value("user_id").(string)
	user_id, err := strconv.Atoi(str_user_id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	data.UserID = user_id
	data.TournamentID = req.TournamentID
	data.TeamName = req.TeamName
	if err := h.service.RequestToJoinTournament(r.Context(), data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Request to join tournament submitted successfully"}, http.StatusOK)
}
