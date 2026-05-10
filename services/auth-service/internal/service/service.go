package service

import (
	"time"

	"auth-service/internal/models"
	pkg_dto "auth-service/internal/pkg"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type TokenService interface {
	Generate(userID, email string) (*pkg_dto.TokenPair, error)
}

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
}

type UserRepository interface {
	Create(user *models.User) error
	GetUser(email string) (*models.User, error)
}

type Service struct {
	hasher     PasswordHasher
	jwt        TokenService
	tokensRepo RefreshTokenRepository
	usersRepo  UserRepository
}

func NewService(hasher PasswordHasher, jwt TokenService, tokensRepo RefreshTokenRepository, usersRepo UserRepository) *Service {
	return &Service{
		hasher:     hasher,
		jwt:        jwt,
		tokensRepo: tokensRepo,
		usersRepo:  usersRepo,
	}
}

func (s *Service) Register(email, password string) (*pkg_dto.TokenResponse, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.usersRepo.Create(user); err != nil {
		return nil, err
	}

	tokens, err := s.jwt.Generate(
		user.ID.String(),
		user.Email,
	)
	if err != nil {
		return nil, err
	}

	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: tokens.ExpiresAt,
		CreatedAt: time.Now().UTC(),
		IsRevoked: false,
	}

	if err := s.tokensRepo.Create(rt); err != nil {
		return nil, err
	}

	return &pkg_dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *Service) Login(email, password string) (*pkg_dto.TokenResponse, error) {
	user, err := s.usersRepo.GetUser(email)
	if err != nil {
		return nil, err
	}

	if flag := s.hasher.Compare(password, user.PasswordHash); flag {
		return nil, err
	}

	tokens, err := s.jwt.Generate(
		user.ID.String(),
		user.Email,
	)

	if err != nil {
		return nil, err
	}

	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: tokens.ExpiresAt,
		CreatedAt: time.Now().UTC(),
		IsRevoked: false,
	}

	if err := s.tokensRepo.Create(rt); err != nil {
		return nil, err
	}

	return &pkg_dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}
