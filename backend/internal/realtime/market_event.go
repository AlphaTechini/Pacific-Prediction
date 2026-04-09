package realtime

import (
	"time"

	"prediction/internal/storage"
)

func NewMarketCreatedEvent(item storage.Market) StreamEvent {
	return StreamEvent{
		Type:       EventTypeMarketCreated,
		OccurredAt: item.CreatedAt,
		MarketID:   item.ID,
		Market:     NewMarketSnapshot(item),
	}
}

func NewMarketUpdatedEvent(item storage.Market, occurredAt time.Time) StreamEvent {
	return StreamEvent{
		Type:       EventTypeMarketUpdated,
		OccurredAt: occurredAt,
		MarketID:   item.ID,
		Market:     NewMarketSnapshot(item),
	}
}

func NewMarketSnapshot(item storage.Market) *MarketSnapshot {
	return &MarketSnapshot{
		ID:                item.ID,
		Title:             item.Title,
		Symbol:            item.Symbol,
		MarketType:        item.MarketType,
		ConditionOperator: item.ConditionOperator,
		ThresholdValue:    item.ThresholdValue,
		SourceType:        item.SourceType,
		SourceInterval:    item.SourceInterval,
		ReferenceValue:    item.ReferenceValue,
		ExpiryTime:        item.ExpiryTime,
		Status:            item.Status,
		Result:            item.Result,
		SettlementValue:   item.SettlementValue,
		ResolvedAt:        item.ResolvedAt,
		ResolutionReason:  item.ResolutionReason,
		CreatedByPlayerID: item.CreatedByPlayerID,
		CreatedAt:         item.CreatedAt,
	}
}
