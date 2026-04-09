package realtime

import (
	"time"

	"prediction/internal/domain"
)

type MarketSnapshot struct {
	ID                domain.MarketID          `json:"id"`
	Title             string                   `json:"title"`
	Symbol            string                   `json:"symbol"`
	MarketType        domain.MarketType        `json:"market_type"`
	ConditionOperator domain.ConditionOperator `json:"condition_operator"`
	ThresholdValue    string                   `json:"threshold_value,omitempty"`
	SourceType        domain.SourceType        `json:"source_type"`
	SourceInterval    string                   `json:"source_interval,omitempty"`
	ReferenceValue    string                   `json:"reference_value,omitempty"`
	ExpiryTime        time.Time                `json:"expiry_time"`
	Status            domain.MarketStatus      `json:"status"`
	Result            domain.MarketResult      `json:"result,omitempty"`
	SettlementValue   string                   `json:"settlement_value,omitempty"`
	ResolvedAt        *time.Time               `json:"resolved_at,omitempty"`
	ResolutionReason  string                   `json:"resolution_reason,omitempty"`
	CreatedByPlayerID domain.PlayerID          `json:"created_by_player_id"`
	CreatedAt         time.Time                `json:"created_at"`
}
