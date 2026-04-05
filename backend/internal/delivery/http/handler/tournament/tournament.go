package tournament

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req domain.TournamentCreateRequest
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	createdBy, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.CreateTournament(r.Context(), createdBy, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Tournament created successfully"}, http.StatusCreated)
}

func (h *TournamentManagerHandler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")
	tournament_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.DeleteTournament(r.Context(), tournament_owner_id, tournament_id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Tournament deleted successfully"}, http.StatusOK)
}

func (h *TournamentManagerHandler) GetTournamentByID(w http.ResponseWriter, r *http.Request) {
	tournament_id := r.URL.Query().Get("id")
	id, err := strconv.Atoi(tournament_id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	tournament, err := h.tournamentService.GetTournamentByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tournament == nil {
		http.Error(w, "Tournament not found", http.StatusNotFound)
		return
	}
	utils.SendData(w, tournament, http.StatusOK)
}

func (h *TournamentManagerHandler) AllTournaments(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	tournaments, err := h.tournamentService.GetAllTournaments(r.Context(), tournament_owner_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, tournaments, http.StatusOK)
}

func (h *TournamentManagerHandler) UpdateTournament(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	var req domain.TournamentCreateRequest
	id := r.URL.Query().Get("id")
	tournament_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid tournament ID", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.UpdateTournament(r.Context(), tournament_owner_id, tournament_id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Tournament updated successfully"}, http.StatusOK)
}
