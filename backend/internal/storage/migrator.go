package storage

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migrator struct {
	pool          *pgxpool.Pool
	migrationsDir string
}

func NewMigrator(pool *pgxpool.Pool, migrationsDir string) *Migrator {
	return &Migrator{
		pool:          pool,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	if err := m.ensureSchemaMigrations(ctx); err != nil {
		return err
	}

	files, err := m.upFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		applied, err := m.isApplied(ctx, file)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		if err := m.applyFile(ctx, file); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) ensureSchemaMigrations(ctx context.Context) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`

	if _, err := m.pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	return nil
}

func (m *Migrator) upFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir(m.migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".up.sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk migrations: %w", err)
	}

	sort.Strings(files)
	return files, nil
}

func (m *Migrator) isApplied(ctx context.Context, path string) (bool, error) {
	version := filepath.Base(path)

	var exists bool
	err := m.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}

	return exists, nil
}

func (m *Migrator) applyFile(ctx context.Context, path string) error {
	version := filepath.Base(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", version, err)
	}

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin migration tx %s: %w", version, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("execute migration %s: %w", version, err)
	}

	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}

	return nil
}
