package auth

import (
	"context"

	"prediction/internal/domain"
)

type contextKey string

const authContextKey contextKey = "auth-context"

func WithAuthContext(ctx context.Context, authContext AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, authContext)
}

func AuthContextFromContext(ctx context.Context) (AuthContext, bool) {
	authContext, ok := ctx.Value(authContextKey).(AuthContext)
	return authContext, ok
}

func RequiredPlayerID(ctx context.Context) (domain.PlayerID, error) {
	authContext, ok := AuthContextFromContext(ctx)
	if !ok {
		return "", ErrUnauthorized
	}

	return authContext.PlayerID, nil
}
