package user

import (
	"context"
	"database/sql"
	"errors"
	"tournament-manager/internal/domain"
	"tournament-manager/internal/domain/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Register(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)
}

func (r *userRepo) GetUserData(ctx context.Context, email, password, role string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 AND role = $2`
	if err := r.db.QueryRowContext(ctx, query, email, role).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, id int, user domain.User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4", user.Username, user.Email, user.Role, id)
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, role FROM users WHERE id = $1", id).Scan(&u.ID, &u.Username, &u.Email, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
