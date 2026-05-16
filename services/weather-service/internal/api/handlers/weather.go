package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) CurrentHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))

	if city == "" {
		http.Error(w, "city is required", http.StatusBadRequest)
		return
	}

	weather, err := h.client.Fetch(city)
	if err == nil {
		w.Header().Set()
	}
}