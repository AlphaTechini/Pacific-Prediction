package position

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"prediction/internal/domain"
)

func NewPositionID() (domain.PositionID, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate position id: %w", err)
	}

	return domain.PositionID("position_" + hex.EncodeToString(buf)), nil
}
