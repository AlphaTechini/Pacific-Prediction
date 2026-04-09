package config

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	Pacifica      PacificaConfig
	Realtime      RealtimeConfig
	Settlement    SettlementConfig
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

type PacificaConfig struct {
	RestBaseURL           string
	MarketInfoCacheTTL    time.Duration
	MarketInfoHTTPTimeout time.Duration
}

type SettlementConfig struct {
	ScanInterval       time.Duration
	ScanBatchSize      int
	PriceLookahead     time.Duration
	PriceRetryInterval time.Duration
}

type RealtimeConfig struct {
	HeartbeatInterval time.Duration
}

func Load() (Config, error) {
	dotEnvMeta, err := loadDotEnv()
	if err != nil {
		return Config{}, err
	}

	authConfig, err := loadAuthConfig()
	if err != nil {
		return Config{}, err
	}

	balanceConfig, err := loadBalanceConfig()
	if err != nil {
		return Config{}, err
	}

	pacificaConfig, err := loadPacificaConfig()
	if err != nil {
		return Config{}, err
	}

	settlementConfig, err := loadSettlementConfig()
	if err != nil {
		return Config{}, err
	}

	realtimeConfig, err := loadRealtimeConfig()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppAddr:       getEnv("APP_ADDR", ":8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		MigrationsDir: resolveMigrationsDir(getEnv("MIGRATIONS_DIR", defaultMigrationsDir()), dotEnvMeta),
		Auth:          authConfig,
		Balance:       balanceConfig,
		Pacifica:      pacificaConfig,
		Realtime:      realtimeConfig,
		Settlement:    settlementConfig,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

type dotEnvLoadResult struct {
	LoadedPaths []string
}

func loadDotEnv() (dotEnvLoadResult, error) {
	result := dotEnvLoadResult{}
	for _, candidate := range dotEnvCandidates() {
		if err := loadDotEnvFile(candidate); err != nil {
			return dotEnvLoadResult{}, err
		}

		if fileExists(candidate) {
			result.LoadedPaths = append(result.LoadedPaths, filepath.Clean(candidate))
		}
	}

	return result, nil
}

func dotEnvCandidates() []string {
	candidates := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("..", "..", ".env"),
	}

	seen := make(map[string]struct{}, len(candidates))
	unique := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		cleaned := filepath.Clean(candidate)
		if _, exists := seen[cleaned]; exists {
			continue
		}

		seen[cleaned] = struct{}{}
		unique = append(unique, cleaned)
	}

	return unique
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("parse %s:%d: expected KEY=VALUE", path, lineNumber)
		}

		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("parse %s:%d: environment key is required", path, lineNumber)
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %s from %s:%d: %w", key, path, lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func defaultMigrationsDir() string {
	candidates := []string{
		"./migrations",
		filepath.Join("..", "..", "migrations"),
	}

	for _, candidate := range candidates {
		if directoryExists(candidate) {
			return candidate
		}
	}

	return "./migrations"
}

func resolveMigrationsDir(value string, dotEnvMeta dotEnvLoadResult) string {
	if value == "" {
		return value
	}

	if filepath.IsAbs(value) {
		return value
	}

	if directoryExists(value) {
		return value
	}

	for _, loadedPath := range dotEnvMeta.LoadedPaths {
		resolved := filepath.Join(filepath.Dir(loadedPath), value)
		if directoryExists(resolved) {
			return resolved
		}
	}

	return value
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
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

func loadPacificaConfig() (PacificaConfig, error) {
	restBaseURL := os.Getenv("PACIFICA_REST_BASE_URL")
	if restBaseURL == "" {
		return PacificaConfig{}, fmt.Errorf("PACIFICA_REST_BASE_URL is required")
	}

	cacheTTLRaw := os.Getenv("PACIFICA_MARKET_INFO_CACHE_TTL")
	if cacheTTLRaw == "" {
		return PacificaConfig{}, fmt.Errorf("PACIFICA_MARKET_INFO_CACHE_TTL is required")
	}

	cacheTTL, err := time.ParseDuration(cacheTTLRaw)
	if err != nil {
		return PacificaConfig{}, fmt.Errorf("parse PACIFICA_MARKET_INFO_CACHE_TTL: %w", err)
	}
	if cacheTTL <= 0 {
		return PacificaConfig{}, fmt.Errorf("PACIFICA_MARKET_INFO_CACHE_TTL must be greater than zero")
	}

	httpTimeoutRaw := os.Getenv("PACIFICA_MARKET_INFO_HTTP_TIMEOUT")
	if httpTimeoutRaw == "" {
		return PacificaConfig{}, fmt.Errorf("PACIFICA_MARKET_INFO_HTTP_TIMEOUT is required")
	}

	httpTimeout, err := time.ParseDuration(httpTimeoutRaw)
	if err != nil {
		return PacificaConfig{}, fmt.Errorf("parse PACIFICA_MARKET_INFO_HTTP_TIMEOUT: %w", err)
	}
	if httpTimeout <= 0 {
		return PacificaConfig{}, fmt.Errorf("PACIFICA_MARKET_INFO_HTTP_TIMEOUT must be greater than zero")
	}

	return PacificaConfig{
		RestBaseURL:           restBaseURL,
		MarketInfoCacheTTL:    cacheTTL,
		MarketInfoHTTPTimeout: httpTimeout,
	}, nil
}

func loadSettlementConfig() (SettlementConfig, error) {
	scanIntervalRaw := os.Getenv("SETTLEMENT_SCAN_INTERVAL")
	if scanIntervalRaw == "" {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_SCAN_INTERVAL is required")
	}

	scanInterval, err := time.ParseDuration(scanIntervalRaw)
	if err != nil {
		return SettlementConfig{}, fmt.Errorf("parse SETTLEMENT_SCAN_INTERVAL: %w", err)
	}
	if scanInterval <= 0 {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_SCAN_INTERVAL must be greater than zero")
	}

	scanBatchSizeRaw := os.Getenv("SETTLEMENT_SCAN_BATCH_SIZE")
	if scanBatchSizeRaw == "" {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_SCAN_BATCH_SIZE is required")
	}

	scanBatchSize, err := strconv.Atoi(scanBatchSizeRaw)
	if err != nil {
		return SettlementConfig{}, fmt.Errorf("parse SETTLEMENT_SCAN_BATCH_SIZE: %w", err)
	}
	if scanBatchSize <= 0 {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_SCAN_BATCH_SIZE must be greater than zero")
	}

	priceLookaheadRaw := os.Getenv("SETTLEMENT_PRICE_LOOKAHEAD")
	if priceLookaheadRaw == "" {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_PRICE_LOOKAHEAD is required")
	}

	priceLookahead, err := time.ParseDuration(priceLookaheadRaw)
	if err != nil {
		return SettlementConfig{}, fmt.Errorf("parse SETTLEMENT_PRICE_LOOKAHEAD: %w", err)
	}
	if priceLookahead <= 0 {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_PRICE_LOOKAHEAD must be greater than zero")
	}

	priceRetryIntervalRaw := os.Getenv("SETTLEMENT_PRICE_RETRY_INTERVAL")
	if priceRetryIntervalRaw == "" {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_PRICE_RETRY_INTERVAL is required")
	}

	priceRetryInterval, err := time.ParseDuration(priceRetryIntervalRaw)
	if err != nil {
		return SettlementConfig{}, fmt.Errorf("parse SETTLEMENT_PRICE_RETRY_INTERVAL: %w", err)
	}
	if priceRetryInterval <= 0 {
		return SettlementConfig{}, fmt.Errorf("SETTLEMENT_PRICE_RETRY_INTERVAL must be greater than zero")
	}

	return SettlementConfig{
		ScanInterval:       scanInterval,
		ScanBatchSize:      scanBatchSize,
		PriceLookahead:     priceLookahead,
		PriceRetryInterval: priceRetryInterval,
	}, nil
}

func loadRealtimeConfig() (RealtimeConfig, error) {
	heartbeatIntervalRaw := os.Getenv("REALTIME_HEARTBEAT_INTERVAL")
	if heartbeatIntervalRaw == "" {
		return RealtimeConfig{}, fmt.Errorf("REALTIME_HEARTBEAT_INTERVAL is required")
	}

	heartbeatInterval, err := time.ParseDuration(heartbeatIntervalRaw)
	if err != nil {
		return RealtimeConfig{}, fmt.Errorf("parse REALTIME_HEARTBEAT_INTERVAL: %w", err)
	}
	if heartbeatInterval <= 0 {
		return RealtimeConfig{}, fmt.Errorf("REALTIME_HEARTBEAT_INTERVAL must be greater than zero")
	}

	return RealtimeConfig{
		HeartbeatInterval: heartbeatInterval,
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
