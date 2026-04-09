package realtime

import "context"

type Publisher interface {
	Publish(ctx context.Context, event StreamEvent) error
}
