package settlement

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func markTemporarySettlementError(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", errSettlementTemporaryFailure, err)
}

func isRetryableSettlementError(err error) bool {
	return errors.Is(err, errSettlementSourceNotReady) || errors.Is(err, errSettlementTemporaryFailure)
}

func resolveWithRetry[T any](
	ctx context.Context,
	retryInterval time.Duration,
	sleep func(context.Context, time.Duration) error,
	resolve func() (T, error),
) (T, error) {
	result, err := resolve()
	if err == nil {
		return result, nil
	}

	var zero T
	if retryInterval <= 0 || !isRetryableSettlementError(err) {
		return zero, err
	}

	if err := sleep(ctx, retryInterval); err != nil {
		return zero, err
	}

	return resolve()
}
