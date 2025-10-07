package user

import (
	"net/http"
	"tournament-manager/utils"
)

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	utils.SendData(w, users, http.StatusOK)
}
