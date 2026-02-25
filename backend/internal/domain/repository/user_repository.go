package repository

import "tournament-manager/internal/domain"

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	Register(user domain.User) error
	Authenticate(email, password, role string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	UpdateUser(id int, user domain.User) error
	DeleteUser(id int) error
}
