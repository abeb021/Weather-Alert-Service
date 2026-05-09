package handlers

import (
	"auth-service/internal/logger"
)

type Handler struct {
	log *logger.Log
}

func NewHandler(log *logger.Log) *Handler {
	return &Handler{
		log: log,
	}
}
