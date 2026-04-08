package position

import (
	"context"

	"prediction/internal/domain"
)

type Controller interface {
	Create(ctx context.Context, playerID domain.PlayerID, input CreateInput) (Record, error)
	ListByPlayerID(ctx context.Context, playerID domain.PlayerID, filter ListFilter) ([]Record, error)
	ValidateCreateInput(ctx context.Context, playerID domain.PlayerID, input CreateInput) error
}
