package models

import (
	"time"
)

type Weather struct {
    City        string    `json:"city"`
    Temperature float64   `json:"temperature"`
    FeelsLike   float64   `json:"feels_like"`
    Condition   string    `json:"condition"`
    WindSpeed   float64   `json:"wind_speed"`
    Humidity    int       `json:"humidity"`
    Timestamp   time.Time `json:"timestamp"`
}