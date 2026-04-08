package pacifica

import "context"

type WebSocketClient interface {
	Connect(ctx context.Context) error
	Subscribe(ctx context.Context, request SubscriptionRequest) error
	ReadEvent(ctx context.Context) (LiveEvent, error)
	SendHeartbeat(ctx context.Context) error
	Close() error
}

type SubscriptionRequest struct {
	Channel  SubscriptionChannel
	Symbol   string
	Interval string
}

type SubscriptionChannel string

const (
	SubscriptionChannelPrices          SubscriptionChannel = "prices"
	SubscriptionChannelCandle          SubscriptionChannel = "candle"
	SubscriptionChannelMarkPriceCandle SubscriptionChannel = "mark_price_candle"
	SubscriptionChannelTrades          SubscriptionChannel = "trades"
)
