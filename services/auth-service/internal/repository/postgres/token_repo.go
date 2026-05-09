package postgres

import (
	"auth-service/internal/models"
	"database/sql"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(tokendbURL string) (*RefreshTokenRepository, error) {
	db, err := sql.Open("postgres", tokendbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &RefreshTokenRepository{db: db}, nil
}

func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	_, err := r.db.Exec(`INSERT INTO refresh_tokens (id, token, expires_at, is_revoked, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		token.ID, token.Token, token.ExpiresAt, token.IsRevoked, token.UserID, token.CreatedAt,
	)
	return err
}
