package settlement

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"prediction/internal/domain"
)

func NewSettlementID() (domain.SettlementID, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate settlement id: %w", err)
	}

	return domain.SettlementID("settlement_" + hex.EncodeToString(buf)), nil
}
