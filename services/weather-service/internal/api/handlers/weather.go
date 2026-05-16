package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	domainErrors "weather-service/internal/domain/errors"
)

func (h *Handler) CurrentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))

	if city == "" {
		writeProblem(w, r, http.StatusBadRequest, "City Required", "The 'city' parameter is missing or empty.")
		return
	}

	weather, err := h.cache.Get(r.Context(), city)
	if err == nil {
		writeJSON(w, http.StatusOK, weather)
		return
	}
	h.log.Debug("cache missing for current weather", "city", city)

	weather, err = h.client.Fetch(city)
	if err != nil {
		if errors.Is(err, domainErrors.ErrCityNotFound) {
			writeProblem(w, r, http.StatusNotFound, "City Not Found", fmt.Sprintf("The city '%s' was not found.", city))
			return
		}
		h.log.Error("failed to fetch current weather", "city", city, "error", err)
		writeProblem(w, r, http.StatusServiceUnavailable, "Service Unavailable", "The weather service is temporarily unavailable.")
		return
	}

	if err := h.cache.Set(r.Context(), city, weather, time.Minute*15); err != nil {
		h.log.Warn("failed to cache current weather", "city", city, "error", err)
	}

	writeJSON(w, http.StatusOK, weather)
}

func (h *Handler) ForecastHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))

	if city == "" {
		writeProblem(w, r, http.StatusBadRequest, "City Required", "The 'city' parameter is missing or empty.")
		return
	}

	weather, err := h.cache.Get(r.Context(), "forecast:"+city)
	if err == nil {
		writeJSON(w, http.StatusOK, weather)
		return
	}
	h.log.Debug("cache missing for forecast", "city", city)

	weather, err = h.client.FetchForecast(city)
	if err != nil {
		if errors.Is(err, domainErrors.ErrCityNotFound) {
			writeProblem(w, r, http.StatusNotFound, "City Not Found", fmt.Sprintf("The city '%s' was not found.", city))
			return
		}
		h.log.Error("failed to fetch forecast", "city", city, "error", err)
		writeProblem(w, r, http.StatusServiceUnavailable, "Service Unavailable", "The weather service is temporarily unavailable.")
		return
	}

	if err := h.cache.Set(r.Context(), "forecast:"+city, weather, time.Minute*30); err != nil {
		h.log.Warn("failed to cache forecast", "city", city, "error", err)
	}

	writeJSON(w, http.StatusOK, weather)
}
