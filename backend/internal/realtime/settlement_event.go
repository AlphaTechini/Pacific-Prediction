package realtime

import "prediction/internal/storage"

func NewMarketSettledEvent(marketItem storage.Market, settlementItem storage.Settlement) StreamEvent {
	resolvedAt := settlementItem.SourceTimestamp
	if marketItem.ResolvedAt != nil {
		resolvedAt = *marketItem.ResolvedAt
	}

	return StreamEvent{
		Type:       EventTypeMarketSettled,
		OccurredAt: settlementItem.CreatedAt,
		MarketID:   marketItem.ID,
		Market:     NewMarketSnapshot(marketItem),
		Settlement: &SettlementSnapshot{
			ID:               settlementItem.ID,
			PacificaSource:   settlementItem.PacificaSource,
			SourceTimestamp:  settlementItem.SourceTimestamp,
			SettlementValue:  settlementItem.SettlementValue,
			Result:           settlementItem.Result,
			ResolutionReason: marketItem.ResolutionReason,
			ResolvedAt:       resolvedAt,
		},
	}
}
