package user

import (
	"encoding/json"
	"net/http"
	"time"
	"tournament-manager/internal/domain"
	"tournament-manager/utils"
)

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	now := time.Now().Format(time.RFC3339)
	user := domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         req.Role,
		CreatedAt:    now,
	}
	if err := h.userService.Register(r.Context(),user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, user, http.StatusCreated)
}
