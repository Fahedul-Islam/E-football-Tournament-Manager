package tournamentmanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

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

	thisRoundDone, err := h.tournamentService.CheckAndAdvanceRound(r.Context(), req.TournamentID, req.Round)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tournament_type, err := h.tournamentService.GetTournamentType(r.Context(), req.TournamentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if (thisRoundDone && tournament_type == "group+knockout") || (thisRoundDone && tournament_type == "knockout") {
		switch req.Round {
		case "Group Stage":
			_, err = h.tournamentService.GenerateKnockoutStage(r.Context(), req.TournamentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case "Round of 16":
			_, err = h.tournamentService.GenerateQuarterFinals(r.Context(), req.TournamentID)
		case "Quarter Finals":
			_, err = h.tournamentService.GenerateSemiFinals(r.Context(), req.TournamentID)
		case "Semifinals":
			_, err = h.tournamentService.GenerateFinal(r.Context(), req.TournamentID)
		case "Final":
			fmt.Println("Tournament has concluded.")
		default:
			http.Error(w, "Unknown round", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	utils.SendData(w, result, http.StatusOK)
}
