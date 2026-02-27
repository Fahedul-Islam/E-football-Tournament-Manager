package user

import (
	"encoding/json"
	"net/http"
	"tournament-manager/utils"
)

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "Email, password and role are required", http.StatusBadRequest)
		return
	}

	loginResponse, err := h.userService.Authenticate(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	utils.SendData(w, loginResponse, http.StatusOK)
}
