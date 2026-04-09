package pacifica

import (
	"errors"
	"fmt"
)

var ErrTemporaryFailure = errors.New("temporary pacifica failure")

func markTemporaryError(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrTemporaryFailure, err)
}

func IsTemporaryError(err error) bool {
	return errors.Is(err, ErrTemporaryFailure)
}
