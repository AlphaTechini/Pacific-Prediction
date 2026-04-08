package player

import (
	"context"
	"fmt"

	"prediction/internal/domain"
	"prediction/internal/storage"
)

type service struct {
	playerRepository storage.PlayerRepository
}

func NewService(playerRepository storage.PlayerRepository) Service {
	return &service{playerRepository: playerRepository}
}

func (s *service) GetProfile(ctx context.Context, playerID domain.PlayerID) (Profile, error) {
	player, err := s.playerRepository.GetByID(ctx, playerID)
	if err != nil {
		return Profile{}, fmt.Errorf("get player profile: %w", err)
	}

	return Profile{
		ID:          player.ID,
		DisplayName: player.DisplayName,
		CreatedAt:   player.CreatedAt,
		UpdatedAt:   player.UpdatedAt,
	}, nil
}
