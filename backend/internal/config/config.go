package config

import (
	"bufio"
	"fmt"
	"math/big"
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
	Database      DatabaseConfig
	MigrationsDir string
	Auth          AuthConfig
	Balance       BalanceConfig
	Market        MarketConfig
	Pacifica      PacificaConfig
	Realtime      RealtimeConfig
	Settlement    SettlementConfig
}

type DatabaseConfig struct {
	URL               string
	MinConns          int32
	MinIdleConns      int32
	MaxConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
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

type MarketConfig struct {
	PriceThresholdCreationBandPercent string
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

	databaseConfig, err := loadDatabaseConfig()
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

	marketConfig, err := loadMarketConfig()
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
		AppAddr:       resolveAppAddr(),
		Database:      databaseConfig,
		MigrationsDir: resolveMigrationsDir(getEnv("MIGRATIONS_DIR", defaultMigrationsDir()), dotEnvMeta),
		Auth:          authConfig,
		Balance:       balanceConfig,
		Market:        marketConfig,
		Pacifica:      pacificaConfig,
		Realtime:      realtimeConfig,
		Settlement:    settlementConfig,
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

func resolveAppAddr() string {
	appAddr := strings.TrimSpace(os.Getenv("APP_ADDR"))
	if appAddr != "" {
		return appAddr
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port != "" {
		return "0.0.0.0:" + port
	}

	return ":3006"
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

func loadDatabaseConfig() (DatabaseConfig, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL is required")
	}

	minConns, err := parseOptionalInt32("DB_MIN_CONNS", 1)
	if err != nil {
		return DatabaseConfig{}, err
	}

	minIdleConns, err := parseOptionalInt32("DB_MIN_IDLE_CONNS", minConns)
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxConns, err := parseOptionalInt32("DB_MAX_CONNS", 4)
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxConnLifetime, err := parseOptionalDuration("DB_MAX_CONN_LIFETIME", 24*time.Hour)
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxConnIdleTime, err := parseOptionalDuration("DB_MAX_CONN_IDLE_TIME", 2*time.Hour)
	if err != nil {
		return DatabaseConfig{}, err
	}

	healthCheckPeriod, err := parseOptionalDuration("DB_HEALTH_CHECK_PERIOD", 30*time.Second)
	if err != nil {
		return DatabaseConfig{}, err
	}

	connectTimeout, err := parseOptionalDuration("DB_CONNECT_TIMEOUT", 15*time.Second)
	if err != nil {
		return DatabaseConfig{}, err
	}

	if minConns < 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_MIN_CONNS must be zero or greater")
	}
	if minIdleConns < 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_MIN_IDLE_CONNS must be zero or greater")
	}
	if maxConns <= 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_MAX_CONNS must be greater than zero")
	}
	if minConns > maxConns {
		return DatabaseConfig{}, fmt.Errorf("DB_MIN_CONNS must be less than or equal to DB_MAX_CONNS")
	}
	if minIdleConns > maxConns {
		return DatabaseConfig{}, fmt.Errorf("DB_MIN_IDLE_CONNS must be less than or equal to DB_MAX_CONNS")
	}
	if maxConnLifetime <= 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_MAX_CONN_LIFETIME must be greater than zero")
	}
	if maxConnIdleTime <= 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_MAX_CONN_IDLE_TIME must be greater than zero")
	}
	if healthCheckPeriod <= 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_HEALTH_CHECK_PERIOD must be greater than zero")
	}
	if connectTimeout <= 0 {
		return DatabaseConfig{}, fmt.Errorf("DB_CONNECT_TIMEOUT must be greater than zero")
	}

	return DatabaseConfig{
		URL:               databaseURL,
		MinConns:          minConns,
		MinIdleConns:      minIdleConns,
		MaxConns:          maxConns,
		MaxConnLifetime:   maxConnLifetime,
		MaxConnIdleTime:   maxConnIdleTime,
		HealthCheckPeriod: healthCheckPeriod,
		ConnectTimeout:    connectTimeout,
	}, nil
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

func parseOptionalInt32(key string, fallback int32) (int32, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return int32(value), nil
}

func parseOptionalDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return value, nil
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

func loadMarketConfig() (MarketConfig, error) {
	bandPercent := strings.TrimSpace(os.Getenv("MARKET_PRICE_THRESHOLD_CREATION_BAND_PERCENT"))
	if bandPercent == "" {
		return MarketConfig{}, fmt.Errorf("MARKET_PRICE_THRESHOLD_CREATION_BAND_PERCENT is required")
	}

	parsed, ok := new(big.Rat).SetString(bandPercent)
	if !ok {
		return MarketConfig{}, fmt.Errorf("MARKET_PRICE_THRESHOLD_CREATION_BAND_PERCENT must be a valid decimal value")
	}

	if parsed.Sign() <= 0 {
		return MarketConfig{}, fmt.Errorf("MARKET_PRICE_THRESHOLD_CREATION_BAND_PERCENT must be greater than zero")
	}

	if parsed.Cmp(big.NewRat(100, 1)) >= 0 {
		return MarketConfig{}, fmt.Errorf("MARKET_PRICE_THRESHOLD_CREATION_BAND_PERCENT must be less than 100")
	}

	return MarketConfig{
		PriceThresholdCreationBandPercent: bandPercent,
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
