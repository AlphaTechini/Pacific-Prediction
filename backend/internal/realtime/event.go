package realtime

import (
	"time"

	"prediction/internal/domain"
)

type StreamEvent struct {
	Type       EventType           `json:"type"`
	OccurredAt time.Time           `json:"occurred_at"`
	MarketID   domain.MarketID     `json:"market_id"`
	Market     *MarketSnapshot     `json:"market,omitempty"`
	Settlement *SettlementSnapshot `json:"settlement,omitempty"`
}
