package user

import (
	"context"
	"errors"
	"tournament-manager/config"
	"tournament-manager/internal/delivery/http/handler/user"
	"tournament-manager/internal/domain"
	"tournament-manager/internal/domain/repository"
	"tournament-manager/utils"
)

type Service interface {
	user.Service
}

// UserService implements the user service layer
type UserRepo interface {
	repository.UserRepository
}
type service struct {
	cfg      *config.Config
	userRepo UserRepo
}

// NewUserService creates a new user service
func NewUserService(cfg *config.Config, userRepo UserRepo) Service {
	return &service{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

// Register registers a new user
func (s *service) Register(ctx context.Context, user domain.User) error {
	if err := utils.IsEmailValid(user.Email); err != nil {
		return err
	}
	if err := utils.ValidatePassword(user.PasswordHash); err != nil {
		return errors.New("Invalid Password Provided")
	}
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return errors.New("Failed to process password")
	}
	user.PasswordHash = hashedPassword
	return s.userRepo.Register(ctx, user)
}

// Authenticate authenticates a user with email, password and role
func (s *service) Authenticate(ctx context.Context, email, password, role string) (*domain.LoginResponse, error) {
	if err:= utils.IsEmailValid(email); err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserData(ctx, email, password, role)
	if err != nil {
		return nil, err
	}

	if err := utils.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, err
	}
	// Generate JWT tokens
	accessToken, refreshToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWT.TokenExpiry.String(),
		TokenType:    "bearer",
	}, nil
}

// GetAllUsers returns all users
func (s *service) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

// GetUserByID returns a user by ID
func (s *service) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// UpdateUser updates a user by ID
func (s *service) UpdateUser(ctx context.Context, id int, user domain.User) error {
	return s.userRepo.UpdateUser(ctx, id, user)
}

// DeleteUser deletes a user by ID
func (s *service) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.DeleteUser(ctx, id)
}
