package tournamentmanager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/domain"
	"tournament-manager/utils"
)

func (h *TournamentManagerHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	str_t_owner_id := r.Context().Value("user_id").(string)
	tournament_owner_id, err := strconv.Atoi(str_t_owner_id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	var req domain.UpadateMatchScoreInput
	err  =json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := h.tournamentService.UpdateScore(tournament_owner_id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	thisRoundDone, err := h.tournamentService.CheckAndAdvanceRound(req.TournamentID, req.Round)
	tournament_type := r.URL.Query().Get("type")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var group_count int
	if req.Round == "Group Stage" {
		group_count, err = h.tournamentService.GetGroupCount(req.TournamentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if (thisRoundDone && tournament_type == "group-knockout") || (thisRoundDone && tournament_type == "knockout") {
		switch req.Round {
			case "Group Stage":
				if group_count == 8	 {
				_, err := h.tournamentService.GenerateRoundOf16(req.TournamentID, 4, nil)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				} else if group_count ==4 {
					_, err := h.tournamentService.GenerateQuarterFinals(req.TournamentID, 2, nil)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else if group_count ==2 {
					_, err := h.tournamentService.GenerateSemiFinals(req.TournamentID, 1, nil)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else if group_count ==1 {
					_, err := h.tournamentService.GenerateFinals(req.TournamentID, 0, nil)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			case "Round of 16":
				_, err := h.tournamentService.GenerateQuarterFinals(req.TournamentID, 2, nil)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "Quarter Finals":
				_, err := h.tournamentService.GenerateSemiFinals(req.TournamentID, 1, nil)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "Semi Finals":
				_, err := h.tournamentService.GenerateFinals(req.TournamentID, 0, nil)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "Final":
				// Tournament is over, no further action needed
			default:
				http.Error(w, "Unknown round", http.StatusBadRequest)
				return
		}
	}
	
	utils.SendData(w, result, http.StatusOK)
}
