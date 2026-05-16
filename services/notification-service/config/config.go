package config

type Config struct {
	Client ClientConfig
	Kafka  KafkaConfig
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type ClientConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func Load() (*Config, error) {
	return nil, nil
}
