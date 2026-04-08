package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"prediction/internal/domain"
)

func NewPlayerID() (domain.PlayerID, error) {
	value, err := newHexID("player")
	if err != nil {
		return "", err
	}

	return domain.PlayerID(value), nil
}

func NewSessionID() (domain.SessionID, error) {
	value, err := newHexID("session")
	if err != nil {
		return "", err
	}

	return domain.SessionID(value), nil
}

func newHexID(prefix string) (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate %s id: %w", prefix, err)
	}

	return prefix + "_" + hex.EncodeToString(buf), nil
}
