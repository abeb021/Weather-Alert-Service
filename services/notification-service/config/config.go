package config

type Config struct {
	Client ClientConfig
}

type ClientConfig struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func Load() (*Config, error) {
	return nil, nil
}
