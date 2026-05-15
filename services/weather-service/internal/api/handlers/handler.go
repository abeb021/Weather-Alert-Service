package handlers

import (
	pkg_dto "weather-service/internal/pkg"
	"log/slog"
)

type Service interface {
	Register(email, password string) (*pkg_dto.TokenResponse, error)
	Login(email, password string) (*pkg_dto.TokenResponse, error)
	Refresh(refreshToken string) (*pkg_dto.TokenResponse, error)
	ValidateAccessToken(accessToken string) (*pkg_dto.ValidateResponse, error)
}

type Handler struct {
	service Service
	log     *slog.Logger
}

func NewHandler(log *slog.Logger, s Service) *Handler {
	return &Handler{
		log:     log,
		service: s,
	}
}
