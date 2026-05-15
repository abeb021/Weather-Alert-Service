package models

import (
	"time"
)

type Weather struct {
	City        string
	Temperature float64
	Condition   string
	WindSpeed   float64
	Humidity    int
	Timestamp   time.Time
}
