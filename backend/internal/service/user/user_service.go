package user

import (
	"tournament-manager/internal/domain"
	"tournament-manager/internal/domain/repository"
	"tournament-manager/internal/delivery/http/handler/user"
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
func (s *service) Register(user domain.User) error {
	return s.userRepo.Register(user)
}

// Authenticate authenticates a user with email, password and role
func (s *service) Authenticate(email, password, role string) (*domain.User, error) {
	return s.userRepo.Authenticate(email, password, role)
}

// GetAllUsers returns all users
func (s *service) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAllUsers()
}

// GetUserByID returns a user by ID
func (s *service) GetUserByID(id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

// UpdateUser updates a user by ID
func (s *service) UpdateUser(id int, user domain.User) error {
	return s.userRepo.UpdateUser(id, user)
}

// DeleteUser deletes a user by ID
func (s *service) DeleteUser(id int) error {
	return s.userRepo.DeleteUser(id)
}
