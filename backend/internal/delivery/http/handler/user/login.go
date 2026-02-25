package user

import (
	"encoding/json"
	"net/http"
	"tournament-manager/utils"
)

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request){
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role	 string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Authenticate(req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.generateToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendData(w, map[string]string{"access_token": accessToken, "refresh_token": refreshToken, "expires_in": h.cfg.JWT.TokenExpiry.String(), "token_type": "bearer"}, http.StatusOK)
}
