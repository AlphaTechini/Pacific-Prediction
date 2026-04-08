package auth

import (
	"strings"

	"prediction/internal/domain"
)

func NormalizeGuestDisplayName(value string, fallbackID domain.PlayerID) string {
	trimmed := strings.TrimSpace(value)
	if trimmed != "" {
		return trimmed
	}

	identifier := string(fallbackID)
	if len(identifier) > 6 {
		identifier = identifier[len(identifier)-6:]
	}

	return "guest-" + identifier
}
