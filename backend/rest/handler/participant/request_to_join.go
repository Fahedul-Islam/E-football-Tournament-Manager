package participant

import (
	"encoding/json"
	"net/http"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *ParticipantHandler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	var req domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.service.RequestToJoinTournament(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Request to join tournament submitted successfully"}, http.StatusOK)
}
