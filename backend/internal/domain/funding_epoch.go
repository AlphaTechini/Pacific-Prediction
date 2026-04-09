package domain

import "time"

func NextFundingEpochFromCreation(createdAt time.Time) time.Time {
	normalizedCreatedAt := NormalizeTime(createdAt)
	if normalizedCreatedAt.IsZero() {
		return time.Time{}
	}

	return normalizedCreatedAt.Truncate(time.Hour).Add(time.Hour)
}
