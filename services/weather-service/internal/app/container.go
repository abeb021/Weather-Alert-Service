package app

import (
	"weather-service/config"
	"weather-service/internal/api/handlers"
	"weather-service/internal/service"
	"log/slog"
	"os"
)

type Container struct {
	Logger     *slog.Logger
	Handler    *handlers.Handler
	JWTService *utils.JWTService
	Cache *postgres.RefreshTokenRepository
}

func NewContainer(logger *slog.Logger, cfg *config.Config) *Container {
	svc := service.NewService(hasher, jwtService, tokensRepo, usersRepo)
	handler := handlers.NewHandler(logger, svc)

	return &Container{
		Logger:     logger,
		Handler:    handler,
	}
}
