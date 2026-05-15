package config

import (
	"time"
)

type Config struct {
	SMTP SMTPConfig
}

type SMTPConfig struct {
	Address           string
	Domain            string
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	MaxMessageBytes   int64
	MaxRecipients     int
	AllowInsecureAuth bool
}

func Load() (*Config, error) {
	return nil, nil
}
