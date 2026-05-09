package app

import (
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/repository/postgres"
)

type Container struct {
	Logger  *logger.Log
	Handler *handlers.Handler
	Repo    *postgres.Repository
}

func NewContainer(logger *logger.Log) *Container {
	handler := handlers.NewHandler(logger)
	repo := postgres.NewRepository()

	return &Container{
		Logger:  logger,
		Handler: handler,
		Repo:    repo,
	}
}
