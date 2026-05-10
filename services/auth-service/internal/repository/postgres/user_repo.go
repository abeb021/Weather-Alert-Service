package postgres

import (
	"auth-service/internal/models"
	"database/sql"
	"errors"
)

var ErrUserCreate = errors.New("failed to create user")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(userdbURL string) (*UserRepository, error) {
	db, err := sql.Open("postgres", userdbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Close() error {
	return r.db.Close()
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Email, user.PasswordHash, user.CreatedAt,
	)
	if err != nil {
		return ErrUserCreate
	}
	return nil
}
