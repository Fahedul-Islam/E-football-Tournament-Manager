package repository

import (
	"context"
	"tournament-manager/internal/domain"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	Register(ctx context.Context, user domain.User) error
	Authenticate(ctx context.Context, email, password, role string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, user domain.User) error
	DeleteUser(ctx context.Context, id int) error
}
