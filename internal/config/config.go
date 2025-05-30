package config

import (
	pc "blockchain-wallet/pkg/db/postgres"
	"os"
)

type Config struct {
	Server struct {
		Address string
		Port    string
	}
	pc.Config
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Устанавливаем адрес сервера по умолчанию
	cfg.Server.Address = "0.0.0.0"
	cfg.Server.Port = "8080"

	// Если есть переменная окружения, используем её
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Address = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}

	return cfg, nil
}
