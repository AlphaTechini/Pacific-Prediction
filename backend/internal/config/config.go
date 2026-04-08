package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv        string
	AppAddr       string
	DatabaseURL   string
	MigrationsDir string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppAddr:       getEnv("APP_ADDR", ":8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "./migrations"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
