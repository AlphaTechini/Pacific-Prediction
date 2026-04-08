package market

import (
	"time"

	"prediction/internal/domain"
)

type Record struct {
	ID                domain.MarketID
	Title             string
	Symbol            string
	MarketType        domain.MarketType
	ConditionOperator domain.ConditionOperator
	ThresholdValue    string
	SourceType        domain.SourceType
	SourceInterval    string
	ReferenceValue    string
	ExpiryTime        time.Time
	Status            domain.MarketStatus
	Result            domain.MarketResult
	SettlementValue   string
	ResolvedAt        *time.Time
	ResolutionReason  string
	CreatedByPlayerID domain.PlayerID
	CreatedAt         time.Time
}

type ListFilter struct {
	Status domain.MarketStatus
	Limit  int
}
