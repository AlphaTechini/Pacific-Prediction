package storage

import (
	"context"
	"fmt"

	"prediction/internal/domain"
)

type SessionPostgresRepository struct {
	queryer Queryer
}

func NewSessionPostgresRepository(queryer Queryer) *SessionPostgresRepository {
	return &SessionPostgresRepository{queryer: queryer}
}

func (r *SessionPostgresRepository) Create(ctx context.Context, input CreateSessionInput) (Session, error) {
	const query = `
INSERT INTO player_sessions (id, player_id, session_token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, player_id, session_token_hash, expires_at, created_at;
`

	var session Session
	if err := r.queryer.QueryRow(ctx, query, string(input.ID), string(input.PlayerID), input.SessionTokenHash, input.ExpiresAt).Scan(
		&session.ID,
		&session.PlayerID,
		&session.SessionTokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
	); err != nil {
		return Session{}, fmt.Errorf("create session: %w", err)
	}

	session.ExpiresAt = domain.NormalizeTime(session.ExpiresAt)
	session.CreatedAt = domain.NormalizeTime(session.CreatedAt)

	return session, nil
}

func (r *SessionPostgresRepository) GetByTokenHash(ctx context.Context, tokenHash string) (Session, error) {
	const query = `
SELECT id, player_id, session_token_hash, expires_at, created_at
FROM player_sessions
WHERE session_token_hash = $1;
`

	var session Session
	if err := r.queryer.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.PlayerID,
		&session.SessionTokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
	); err != nil {
		return Session{}, fmt.Errorf("get session by token hash: %w", err)
	}

	session.ExpiresAt = domain.NormalizeTime(session.ExpiresAt)
	session.CreatedAt = domain.NormalizeTime(session.CreatedAt)

	return session, nil
}
