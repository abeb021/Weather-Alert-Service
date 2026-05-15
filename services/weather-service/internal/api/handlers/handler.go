package handlers

import (
	"log/slog"
)

type Handler struct {
	log     *slog.Logger
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		log:     log,
	}
}
