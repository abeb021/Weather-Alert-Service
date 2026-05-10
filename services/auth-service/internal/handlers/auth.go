package handlers

import (
	pkg_dto "auth-service/internal/pkg"
	"encoding/json"
	"net/http"
)

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req pkg_dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		h.log.Logger.Error("register request body decoding error", "error", err)
		return
	}

	token, err := h.service.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Logger.Error("service error", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(token)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req pkg_dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		h.log.Logger.Error("login request body decoding error", "error", err)
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Logger.Error("service error", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(token)
}
