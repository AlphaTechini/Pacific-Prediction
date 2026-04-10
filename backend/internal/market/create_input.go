package market

import (
	"time"

	"prediction/internal/domain"
)

type CreateInput struct {
	Title                string
	Symbol               string
	SymbolPriceIncrement string
	MarketType           domain.MarketType
	ConditionOperator    domain.ConditionOperator
	CreatorSide          domain.PositionSide
	CreatorStakeAmount   string
	ThresholdValue       string
	SourceType           domain.SourceType
	SourceInterval       string
	ReferenceValue       string
	ExpiryTime           time.Time
	CreatedByPlayerID    domain.PlayerID
}
