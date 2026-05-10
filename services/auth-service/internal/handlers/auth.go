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
		return
	}

	token, err := h.service.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(token)
}
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}
