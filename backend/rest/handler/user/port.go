package user

import (
	"tournament-manager/config"
	"tournament-manager/domain"
)

type Service interface {
	Register(user domain.User) error
	Authenticate(email, password, role string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	UpdateUser(id int, user domain.User) error
	DeleteUser(id int) error
}

type UserHandler struct {
	cfg         *config.Config
	userService Service
}

func NewUserHandler(cfg *config.Config, userService Service) *UserHandler {
	return &UserHandler{
		cfg:         cfg,
		userService: userService,
	}
}
