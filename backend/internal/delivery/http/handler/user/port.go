package user

import (
	"context"
	"tournament-manager/internal/domain"
)

type Service interface {
	Register(ctx context.Context, user domain.User) error
	Authenticate(ctx context.Context, email, password, role string) (*domain.LoginResponse, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, user domain.User) error
	DeleteUser(ctx context.Context, id int) error
}

type UserHandler struct {
	userService Service
}

func NewUserHandler(userService Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
