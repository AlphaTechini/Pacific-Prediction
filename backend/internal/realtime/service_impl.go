package realtime

import (
	"context"
	"sync"
)

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Subscribe(ctx context.Context) (Subscription, error) {
	events := make(chan StreamEvent)
	done := make(chan struct{})

	var closeOnce sync.Once
	closeSubscription := func() {
		closeOnce.Do(func() {
			close(done)
		})
	}

	go func() {
		defer close(events)

		select {
		case <-ctx.Done():
		case <-done:
		}
	}()

	return Subscription{
		Events: events,
		Close:  closeSubscription,
	}, nil
}
