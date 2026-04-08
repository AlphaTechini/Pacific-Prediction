package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid domain input")
	ErrInvalidState = errors.New("invalid domain state")
	ErrNotFound     = errors.New("domain record not found")
	ErrConflict     = errors.New("domain conflict")
)

type ValidationError struct {
	Field   string
	Message string
	Value   any
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e ValidationError) Unwrap() error {
	return ErrInvalidInput
}

func NewValidationError(field, message string, value any) error {
	return ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}
