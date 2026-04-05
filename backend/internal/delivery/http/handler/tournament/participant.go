package tournament

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) GetAllParticipant(w http.ResponseWriter, r *http.Request) {
	str_id := r.URL.Query().Get("tournament_id")
	tournament_id, err := strconv.Atoi(str_id)
	if err != nil {
		http.Error(w, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	participant, err := h.tournamentService.GetAllParticipant(r.Context(), tournament_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendData(w, participant, http.StatusOK)
}

func (h *TournamentManagerHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tournament_owner_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	var participant domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.AddParticipant(r.Context(), tournament_owner_id, participant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant added successfully"}, http.StatusOK)
}

func (h *TournamentManagerHandler) ApproveParticipant(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tournament_owner_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var req domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.ApproveParticipant(r.Context(), tournament_owner_id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant approved successfully"}, http.StatusOK)
}

func (h *TournamentManagerHandler) RejectParticipant(w http.ResponseWriter, r *http.Request) {
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
	var req domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.RejectParticipant(r.Context(), tournament_owner_id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant rejected successfully"}, http.StatusOK)
}

func (h *TournamentManagerHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
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
	var req domain.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.tournamentService.RemoveParticipant(r.Context(), tournament_owner_id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Participant removed successfully"}, http.StatusOK)
}
