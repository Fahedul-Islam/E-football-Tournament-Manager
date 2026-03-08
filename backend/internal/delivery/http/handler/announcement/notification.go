package announcement

import (
	"net/http"
	"strconv"
	"tournament-manager/utils"
)

func (h *AnnouncementHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
		page = p
	}

	notifications, err := h.announcementService.GetNotifications(r.Context(), userID, page)
	if err != nil {
		http.Error(w, "Failed to get notifications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, notifications, http.StatusOK)
}

func (h *AnnouncementHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	notificationIDStr := r.URL.Query().Get("notification_id")
	if notificationIDStr == "" {
		http.Error(w, "Missing notification_id", http.StatusBadRequest)
		return
	}

	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		http.Error(w, "Invalid notification_id", http.StatusBadRequest)
		return
	}

	err = h.announcementService.MarkNotificationAsRead(r.Context(), notificationID, userID)
	if err != nil {
		http.Error(w, "Failed to mark notification as read: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnnouncementHandler) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	err = h.announcementService.MarkAllNotificationsAsRead(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to mark all notifications as read: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
