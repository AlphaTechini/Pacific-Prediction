package pacifica

import (
	"context"
	"fmt"
	"time"

	"prediction/internal/domain"
)

type subscriptionManager struct {
	clientFactory        WebSocketClientFactory
	heartbeatInterval    time.Duration
	reconnectDelay       time.Duration
	maxReconnectAttempts int
}

func (m *subscriptionManager) Run(ctx context.Context, requests []SubscriptionRequest, handler LiveEventHandler) error {
	if len(requests) == 0 {
		return domain.NewValidationError("requests", "at least one subscription is required", requests)
	}
	if handler == nil {
		return domain.NewValidationError("handler", "live event handler is required", nil)
	}
	if m.clientFactory == nil {
		return domain.NewValidationError("client_factory", "websocket client factory is required", nil)
	}

	attempts := 0
	for {
		client := m.clientFactory.NewWebSocketClient()
		if client == nil {
			return domain.NewValidationError("client_factory", "websocket client factory returned nil client", nil)
		}

		err := m.runSession(ctx, client, requests, handler)
		_ = client.Close()

		if ctx.Err() != nil {
			return nil
		}
		if err == nil {
			return nil
		}
		if !m.shouldReconnect(attempts) {
			return fmt.Errorf("subscription manager stopped after reconnect attempts: %w", err)
		}

		attempts++
		timer := time.NewTimer(m.reconnectDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
}

func (m *subscriptionManager) runSession(ctx context.Context, client WebSocketClient, requests []SubscriptionRequest, handler LiveEventHandler) error {
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("connect websocket client: %w", err)
	}

	for _, request := range requests {
		if err := client.Subscribe(ctx, request); err != nil {
			return fmt.Errorf("subscribe websocket client: %w", err)
		}
	}

	heartbeatTicker := time.NewTicker(m.heartbeatInterval)
	defer heartbeatTicker.Stop()

	eventCh := make(chan LiveEvent)
	errCh := make(chan error, 1)

	go m.readLoop(ctx, client, eventCh, errCh)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if err == nil {
				return nil
			}
			return err
		case event := <-eventCh:
			if err := handler(ctx, event); err != nil {
				return fmt.Errorf("handle live event: %w", err)
			}
		case <-heartbeatTicker.C:
			if err := client.SendHeartbeat(ctx); err != nil {
				return fmt.Errorf("send websocket heartbeat: %w", err)
			}
		}
	}
}

func (m *subscriptionManager) readLoop(ctx context.Context, client WebSocketClient, eventCh chan<- LiveEvent, errCh chan<- error) {
	for {
		event, err := client.ReadEvent(ctx)
		if err != nil {
			select {
			case errCh <- fmt.Errorf("read websocket event: %w", err):
			default:
			}
			return
		}

		select {
		case <-ctx.Done():
			return
		case eventCh <- event:
		}
	}
}

func (m *subscriptionManager) shouldReconnect(attempts int) bool {
	if m.maxReconnectAttempts <= 0 {
		return true
	}

	return attempts < m.maxReconnectAttempts
}
