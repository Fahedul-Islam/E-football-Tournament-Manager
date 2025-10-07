package tournamentmanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req domain.TournamentCreateRequest
	// Extract the created_by field from the request context or token
	userIDStr := r.Context().Value("user_id").(string)
	createdBy, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("CreateTournament request:", req)
	if err := h.tournamentService.CreateTournament(createdBy, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Tournament created successfully"}, http.StatusCreated)
}
