package realtime

import (
	"context"
	"sync"
	"time"
)

const defaultSubscriberBufferSize = 32

type hub struct {
	mu          sync.RWMutex
	nextID      int64
	subscribers map[int64]chan StreamEvent
}

type subscriberSnapshot struct {
	id int64
	ch chan StreamEvent
}

func NewHub() Hub {
	return &hub{
		subscribers: make(map[int64]chan StreamEvent),
	}
}

func NewService() Service {
	return NewHub()
}

func (h *hub) Subscribe(ctx context.Context) (Subscription, error) {
	ch := make(chan StreamEvent, defaultSubscriberBufferSize)
	subscriberID := h.addSubscriber(ch)

	var closeOnce sync.Once
	closeSubscription := func() {
		closeOnce.Do(func() {
			h.removeSubscriber(subscriberID)
		})
	}

	go func() {
		<-ctx.Done()
		closeSubscription()
	}()

	return Subscription{
		Events: ch,
		Close:  closeSubscription,
	}, nil
}

func (h *hub) Publish(ctx context.Context, event StreamEvent) error {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}

	h.mu.RLock()
	snapshots := make([]subscriberSnapshot, 0, len(h.subscribers))
	for id, ch := range h.subscribers {
		snapshots = append(snapshots, subscriberSnapshot{
			id: id,
			ch: ch,
		})
	}
	h.mu.RUnlock()

	for _, subscriber := range snapshots {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case subscriber.ch <- event:
		default:
			h.removeSubscriber(subscriber.id)
		}
	}

	return nil
}

func (h *hub) addSubscriber(ch chan StreamEvent) int64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.nextID++
	subscriberID := h.nextID
	h.subscribers[subscriberID] = ch
	return subscriberID
}

func (h *hub) removeSubscriber(id int64) {
	h.mu.Lock()
	ch, ok := h.subscribers[id]
	if ok {
		delete(h.subscribers, id)
	}
	h.mu.Unlock()

	if ok {
		close(ch)
	}
}
