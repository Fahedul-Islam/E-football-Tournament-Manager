package userrepo

import (
	"database/sql"
	"tournament-manager/domain"
	"tournament-manager/rest/handler/user"
	"tournament-manager/utils"
)

type UserRepo interface {
	user.Service
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Register(user domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)
}

func (r *userRepo) Authenticate(email, password, role string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 AND role = $2`
	if err := r.db.QueryRow(query, email, role).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		return nil, err
	}

	if err := utils.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdateUser(id int, user domain.User) error {
	_, err := r.db.Exec("UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4", user.Username, user.Email, user.Role, id)
	return err
}

func (r *userRepo) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepo) GetAllUsers() ([]*domain.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *userRepo) GetUserByID(id int) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow("SELECT id, username, email, role FROM users WHERE id = $1", id).Scan(&u.ID, &u.Username, &u.Email, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
