package domain

import "time"

func NormalizeTime(value time.Time) time.Time {
	if value.IsZero() {
		return value
	}

	return value.UTC().Round(0)
}

func NowUTC() time.Time {
	return NormalizeTime(time.Now())
}
