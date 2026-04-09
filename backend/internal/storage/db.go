package storage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"prediction/internal/config"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, databaseConfig config.DatabaseConfig) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	poolConfig.MinConns = databaseConfig.MinConns
	poolConfig.MinIdleConns = databaseConfig.MinIdleConns
	poolConfig.MaxConns = databaseConfig.MaxConns
	poolConfig.MaxConnLifetime = databaseConfig.MaxConnLifetime
	poolConfig.MaxConnIdleTime = databaseConfig.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = databaseConfig.HealthCheckPeriod
	poolConfig.ConnConfig.ConnectTimeout = databaseConfig.ConnectTimeout

	endpoint := databaseEndpointLabel(databaseConfig.URL)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("open database pool for %s: %w", endpoint, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database %s: %w", endpoint, err)
	}

	return &DB{pool: pool}, nil
}

func databaseEndpointLabel(databaseURL string) string {
	parsed, err := url.Parse(databaseURL)
	if err != nil || parsed.Host == "" {
		return "configured database"
	}

	return parsed.Host
}

func (db *DB) Close() {
	if db == nil || db.pool == nil {
		return
	}

	db.pool.Close()
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}
