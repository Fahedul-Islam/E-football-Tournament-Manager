package announcement

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *AnnouncementHandler) CommentOnAnnouncement(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
	if err != nil {
		http.Error(w, "Invalid announcement_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	var req domain.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.announcementService.CreateComment(r.Context(), tournamentID, announcementID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, comment, http.StatusOK)
}

func (h *AnnouncementHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	announcementID, err := strconv.Atoi(r.URL.Query().Get("announcement_id"))
	if err != nil {
		http.Error(w, "Invalid announcement_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	parentCommentIDStr := r.URL.Query().Get("parent_comment_id")
	var parentCommentID *int
	if parentCommentIDStr != "" {
		id, err := strconv.Atoi(parentCommentIDStr)
		if err != nil {
			http.Error(w, "Invalid parent_comment_id", http.StatusBadRequest)
			return
		}
		parentCommentID = &id
	}

	comments, err := h.announcementService.GetComments(r.Context(), tournamentID, userID, parentCommentID, announcementID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, comments, http.StatusOK)
}

func (h *AnnouncementHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.announcementService.DeleteComment(r.Context(), tournamentID, userID, commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"message": "Comment deleted successfully"}, http.StatusOK)
}

func (h *AnnouncementHandler) EditComment(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	var req domain.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.announcementService.EditComment(r.Context(), tournamentID, userID, commentID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, comment, http.StatusOK)
}

func (h *AnnouncementHandler) ReactToComment(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := strconv.Atoi(r.URL.Query().Get("tournament_id"))
	if err != nil {
		http.Error(w, "Invalid tournament_id", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value("user_id").(string))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	reaction := r.URL.Query().Get("reaction")
	if reaction == "" {
		http.Error(w, "Reaction is required", http.StatusBadRequest)
		return
	}

	comment, err := h.announcementService.ReactToComment(r.Context(), tournamentID, commentID, userID, reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, comment, http.StatusOK)
}

