package config

import (
	"os"
)

type Config struct {
	Server struct {
		Address string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Устанавливаем адрес сервера по умолчанию
	cfg.Server.Address = ":8080"

	// Если есть переменная окружения, используем её
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Address = ":" + port
	}

	return cfg, nil
}
