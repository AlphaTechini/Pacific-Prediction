package realtime

import "context"

type Controller interface {
	Subscribe(ctx context.Context) (Subscription, error)
}
