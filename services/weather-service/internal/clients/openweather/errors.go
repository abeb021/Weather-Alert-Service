package openweather

import "errors"

var (
	ErrCityNotFound   = errors.New("city not found")
	ErrInvalidAPIKey  = errors.New("invalid API key")
	ErrAPIUnavailable = errors.New("weather API unavailable")
)
