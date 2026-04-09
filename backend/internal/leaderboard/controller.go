package leaderboard

import "context"

type Controller interface {
	GetSnapshot(ctx context.Context, limit int) (Snapshot, error)
}
