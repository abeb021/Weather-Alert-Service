package errors

import "errors"

var (
    ErrCityRequired = errors.New("city is required")
    ErrCityNotFound   = errors.New("city not found")
	ErrInvalidAPIKey  = errors.New("invalid API key")
	ErrAPIUnavailable = errors.New("weather API unavailable")
)
