package storage

import (
	"context"
	"fmt"

	"prediction/internal/domain"
)

type PlayerPostgresRepository struct {
	queryer Queryer
}

func NewPlayerPostgresRepository(queryer Queryer) *PlayerPostgresRepository {
	return &PlayerPostgresRepository{queryer: queryer}
}

func (r *PlayerPostgresRepository) Create(ctx context.Context, input CreatePlayerInput) (Player, error) {
	const query = `
INSERT INTO players (id, display_name)
VALUES ($1, $2)
RETURNING id, display_name, created_at, updated_at;
`

	var player Player
	if err := r.queryer.QueryRow(ctx, query, string(input.ID), input.DisplayName).Scan(
		&player.ID,
		&player.DisplayName,
		&player.CreatedAt,
		&player.UpdatedAt,
	); err != nil {
		return Player{}, fmt.Errorf("create player: %w", err)
	}

	player.CreatedAt = domain.NormalizeTime(player.CreatedAt)
	player.UpdatedAt = domain.NormalizeTime(player.UpdatedAt)

	return player, nil
}

func (r *PlayerPostgresRepository) GetByID(ctx context.Context, playerID domain.PlayerID) (Player, error) {
	const query = `
SELECT id, display_name, created_at, updated_at
FROM players
WHERE id = $1;
`

	var player Player
	if err := r.queryer.QueryRow(ctx, query, string(playerID)).Scan(
		&player.ID,
		&player.DisplayName,
		&player.CreatedAt,
		&player.UpdatedAt,
	); err != nil {
		return Player{}, fmt.Errorf("get player by id: %w", err)
	}

	player.CreatedAt = domain.NormalizeTime(player.CreatedAt)
	player.UpdatedAt = domain.NormalizeTime(player.UpdatedAt)

	return player, nil
}
