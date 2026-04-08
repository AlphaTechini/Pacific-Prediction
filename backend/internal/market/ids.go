package market

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"prediction/internal/domain"
)

func NewMarketID() (domain.MarketID, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate market id: %w", err)
	}

	return domain.MarketID("market_" + hex.EncodeToString(buf)), nil
}
