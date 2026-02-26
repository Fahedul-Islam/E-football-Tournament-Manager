package user

import (
	"context"
	"tournament-manager/config"
	"tournament-manager/internal/domain"
)

type Service interface {
	Register(ctx context.Context,user domain.User) error
	Authenticate(ctx context.Context,email, password, role string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context,id int) (*domain.User, error)
	UpdateUser(ctx context.Context,id int, user domain.User) error
	DeleteUser(ctx context.Context,id int) error
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
