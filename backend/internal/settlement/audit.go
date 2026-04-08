package settlement

import (
	"time"

	"prediction/internal/domain"
)

type AuditRecord struct {
	MarketID        domain.MarketID
	PacificaSource  string
	SourceTimestamp time.Time
	RawPayload      []byte
	SettlementValue string
	Result          domain.MarketResult
}
