package domain

import (
	"fmt"
	"math/big"
	"strings"
)

func ParseDecimal(value string) (*big.Rat, error) {
	parsed, ok := new(big.Rat).SetString(strings.TrimSpace(value))
	if !ok {
		return nil, fmt.Errorf("parse decimal %q", value)
	}

	return parsed, nil
}

func FitsNumericScale(value string, maxScale int) bool {
	parts := strings.SplitN(strings.TrimSpace(value), ".", 2)
	if len(parts) < 2 {
		return true
	}

	return len(parts[1]) <= maxScale
}

func FormatFixedScaleDecimal(value *big.Rat, scale int) string {
	if value == nil {
		return ""
	}

	return value.FloatString(scale)
}

func CalculateEvenOddsPayout(stakeAmount string) (string, error) {
	stake, err := ParseDecimal(stakeAmount)
	if err != nil {
		return "", err
	}

	payout := new(big.Rat).Mul(stake, big.NewRat(2, 1))
	return FormatFixedScaleDecimal(payout, 8), nil
}
