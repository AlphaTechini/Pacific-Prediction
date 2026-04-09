package domain

import (
	"fmt"
	"strings"
	"time"
)

func CandleIntervalDuration(interval string) (time.Duration, error) {
	switch strings.TrimSpace(interval) {
	case "1m":
		return time.Minute, nil
	case "3m":
		return 3 * time.Minute, nil
	case "5m":
		return 5 * time.Minute, nil
	case "15m":
		return 15 * time.Minute, nil
	case "30m":
		return 30 * time.Minute, nil
	case "1h":
		return time.Hour, nil
	case "2h":
		return 2 * time.Hour, nil
	case "4h":
		return 4 * time.Hour, nil
	case "8h":
		return 8 * time.Hour, nil
	case "12h":
		return 12 * time.Hour, nil
	case "1d":
		return 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported candle interval %q", interval)
	}
}

func CandleWindowForExpiry(expiry time.Time, interval string) (time.Time, time.Time, error) {
	duration, err := CandleIntervalDuration(interval)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	normalizedExpiry := NormalizeTime(expiry)
	if normalizedExpiry.IsZero() {
		return time.Time{}, time.Time{}, fmt.Errorf("expiry time is required")
	}

	if normalizedExpiry.UnixNano()%duration.Nanoseconds() != 0 {
		return time.Time{}, time.Time{}, fmt.Errorf(
			"expiry time %s is not aligned to candle interval %s",
			normalizedExpiry.Format(time.RFC3339Nano),
			interval,
		)
	}

	return normalizedExpiry.Add(-duration), normalizedExpiry, nil
}

func NextCandleExpiryFromCreation(createdAt time.Time, interval string) (time.Time, error) {
	duration, err := CandleIntervalDuration(interval)
	if err != nil {
		return time.Time{}, err
	}

	normalizedCreatedAt := NormalizeTime(createdAt)
	if normalizedCreatedAt.IsZero() {
		return time.Time{}, fmt.Errorf("created time is required")
	}

	target := normalizedCreatedAt.Add(duration)
	remainder := target.UnixNano() % duration.Nanoseconds()
	if remainder == 0 {
		return target, nil
	}

	return target.Add(time.Duration(duration.Nanoseconds() - remainder)), nil
}
