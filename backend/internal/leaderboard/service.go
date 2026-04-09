package leaderboard

import "context"

type Service interface {
	GetSnapshot(ctx context.Context, limit int) (Snapshot, error)
}
