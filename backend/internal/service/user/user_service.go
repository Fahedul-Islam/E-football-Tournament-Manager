package user

import (
	"context"
	"tournament-manager/internal/delivery/http/handler/user"
	"tournament-manager/internal/domain"
	"tournament-manager/internal/domain/repository"
)

type Service interface {
	user.Service
}

// UserService implements the user service layer
type UserRepo interface {
	repository.UserRepository
}
type service struct {
	userRepo UserRepo
}

// NewUserService creates a new user service
func NewUserService(userRepo UserRepo) Service {
	return &service{
		userRepo: userRepo,
	}
}

// Register registers a new user
func (s *service) Register(ctx context.Context, user domain.User) error {
	return s.userRepo.Register(ctx, user)
}

// Authenticate authenticates a user with email, password and role
func (s *service) Authenticate(ctx context.Context,email, password, role string) (*domain.User, error) {
	return s.userRepo.Authenticate(ctx,email, password, role)
}

// GetAllUsers returns all users
func (s *service) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

// GetUserByID returns a user by ID
func (s *service) GetUserByID(ctx context.Context,id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx,id)
}

// UpdateUser updates a user by ID
func (s *service) UpdateUser(ctx context.Context,id int, user domain.User) error {
	return s.userRepo.UpdateUser(ctx,id, user)
}

// DeleteUser deletes a user by ID
func (s *service) DeleteUser(ctx context.Context,id int) error {
	return s.userRepo.DeleteUser(ctx,id)
}
