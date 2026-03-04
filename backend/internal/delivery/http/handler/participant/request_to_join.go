package participant

import (
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *ParticipantHandler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	var data domain.ParticipantRequest
	tournament_id := r.URL.Query().Get("tournament_id")
	tournamentID, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	team_name := r.URL.Query().Get("team_name")
	if team_name == "" {
		http.Error(w, "Team name is required", http.StatusBadRequest)
		return
	}
	str_user_id := r.Context().Value("user_id").(string)
	user_id, err := strconv.Atoi(str_user_id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	data.UserID = user_id
	data.TournamentID = tournamentID
	data.TeamName = team_name
	if err := h.service.RequestToJoinTournament(r.Context(), data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Request to join tournament submitted successfully"}, http.StatusOK)
}
