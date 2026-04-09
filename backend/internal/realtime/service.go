package realtime

import "context"

type Service interface {
	Subscribe(ctx context.Context) (Subscription, error)
}
