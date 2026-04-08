package pacifica

import (
	"context"
	"time"
)

type SubscriptionManager interface {
	Run(ctx context.Context, requests []SubscriptionRequest, handler LiveEventHandler) error
}

type LiveEventHandler func(ctx context.Context, event LiveEvent) error

type WebSocketClientFactory interface {
	NewWebSocketClient() WebSocketClient
}

type SubscriptionManagerConfig struct {
	HeartbeatInterval    time.Duration
	ReconnectDelay       time.Duration
	MaxReconnectAttempts int
}

func DefaultSubscriptionManagerConfig() SubscriptionManagerConfig {
	return SubscriptionManagerConfig{
		HeartbeatInterval:    30 * time.Second,
		ReconnectDelay:       3 * time.Second,
		MaxReconnectAttempts: 0,
	}
}

func NewSubscriptionManager(factory WebSocketClientFactory, config SubscriptionManagerConfig) SubscriptionManager {
	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = 30 * time.Second
	}
	if config.ReconnectDelay <= 0 {
		config.ReconnectDelay = 3 * time.Second
	}

	return &subscriptionManager{
		clientFactory:        factory,
		heartbeatInterval:    config.HeartbeatInterval,
		reconnectDelay:       config.ReconnectDelay,
		maxReconnectAttempts: config.MaxReconnectAttempts,
	}
}
