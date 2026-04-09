package realtime

import (
	"time"

	"prediction/internal/domain"
)

type SettlementSnapshot struct {
	ID               domain.SettlementID `json:"id"`
	PacificaSource   string              `json:"pacifica_source"`
	SourceTimestamp  time.Time           `json:"source_timestamp"`
	SettlementValue  string              `json:"settlement_value,omitempty"`
	Result           domain.MarketResult `json:"result"`
	ResolutionReason string              `json:"resolution_reason,omitempty"`
	ResolvedAt       time.Time           `json:"resolved_at"`
}
