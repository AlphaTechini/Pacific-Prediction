package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv        string
	AppAddr       string
	DatabaseURL   string
	MigrationsDir string
	Auth          AuthConfig
	Balance       BalanceConfig
}

type AuthConfig struct {
	SessionTTL     time.Duration
	CookieName     string
	CookieSecure   bool
	CookieSameSite http.SameSite
	CookieDomain   string
	CookiePath     string
}

type BalanceConfig struct {
	StartingBalance string
}

func Load() (Config, error) {
	authConfig, err := loadAuthConfig()
	if err != nil {
		return Config{}, err
	}

	balanceConfig, err := loadBalanceConfig()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppAddr:       getEnv("APP_ADDR", ":8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "./migrations"),
		Auth:          authConfig,
		Balance:       balanceConfig,
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

func loadAuthConfig() (AuthConfig, error) {
	sessionTTLRaw := os.Getenv("AUTH_SESSION_TTL")
	if sessionTTLRaw == "" {
		return AuthConfig{}, fmt.Errorf("AUTH_SESSION_TTL is required")
	}

	sessionTTL, err := time.ParseDuration(sessionTTLRaw)
	if err != nil {
		return AuthConfig{}, fmt.Errorf("parse AUTH_SESSION_TTL: %w", err)
	}
	if sessionTTL <= 0 {
		return AuthConfig{}, fmt.Errorf("AUTH_SESSION_TTL must be greater than zero")
	}

	cookieName := os.Getenv("AUTH_COOKIE_NAME")
	if cookieName == "" {
		return AuthConfig{}, fmt.Errorf("AUTH_COOKIE_NAME is required")
	}

	cookieSecureRaw := os.Getenv("AUTH_COOKIE_SECURE")
	if cookieSecureRaw == "" {
		return AuthConfig{}, fmt.Errorf("AUTH_COOKIE_SECURE is required")
	}

	cookieSecure, err := strconv.ParseBool(cookieSecureRaw)
	if err != nil {
		return AuthConfig{}, fmt.Errorf("parse AUTH_COOKIE_SECURE: %w", err)
	}

	cookieSameSite, err := parseSameSite(os.Getenv("AUTH_COOKIE_SAME_SITE"))
	if err != nil {
		return AuthConfig{}, err
	}

	cookiePath := os.Getenv("AUTH_COOKIE_PATH")
	if cookiePath == "" {
		return AuthConfig{}, fmt.Errorf("AUTH_COOKIE_PATH is required")
	}

	return AuthConfig{
		SessionTTL:     sessionTTL,
		CookieName:     cookieName,
		CookieSecure:   cookieSecure,
		CookieSameSite: cookieSameSite,
		CookieDomain:   os.Getenv("AUTH_COOKIE_DOMAIN"),
		CookiePath:     cookiePath,
	}, nil
}

func loadBalanceConfig() (BalanceConfig, error) {
	startingBalance := os.Getenv("PLAYER_STARTING_BALANCE")
	if startingBalance == "" {
		return BalanceConfig{}, fmt.Errorf("PLAYER_STARTING_BALANCE is required")
	}

	return BalanceConfig{
		StartingBalance: startingBalance,
	}, nil
}

func parseSameSite(value string) (http.SameSite, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "lax":
		return http.SameSiteLaxMode, nil
	case "strict":
		return http.SameSiteStrictMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return http.SameSiteDefaultMode, fmt.Errorf("AUTH_COOKIE_SAME_SITE must be one of lax, strict, or none")
	}
}
