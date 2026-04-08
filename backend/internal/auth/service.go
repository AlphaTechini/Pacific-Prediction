package auth

import (
	"context"
	"fmt"

	"prediction/internal/config"
	"prediction/internal/domain"
	"prediction/internal/storage"

	"github.com/jackc/pgx/v5"
)

type ServiceDeps struct {
	Config            config.Config
	TxManager         *storage.TxManager
	SessionRepository storage.SessionRepository
}

type service struct {
	config            config.Config
	txManager         *storage.TxManager
	sessionRepository storage.SessionRepository
}

func NewService(deps ServiceDeps) Service {
	return &service{
		config:            deps.Config,
		txManager:         deps.TxManager,
		sessionRepository: deps.SessionRepository,
	}
}

func (s *service) CreateGuestSession(ctx context.Context, input CreateGuestSessionInput) (IssuedSession, error) {
	playerID, err := NewPlayerID()
	if err != nil {
		return IssuedSession{}, err
	}

	sessionID, err := NewSessionID()
	if err != nil {
		return IssuedSession{}, err
	}

	rawToken, err := GenerateOpaqueToken()
	if err != nil {
		return IssuedSession{}, err
	}

	expiresAt := domain.NowUTC().Add(s.config.Auth.SessionTTL)
	displayName := NormalizeGuestDisplayName(input.DisplayName, playerID)
	sessionTokenHash := HashToken(rawToken)

	var createdSession storage.Session
	if err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		playerRepository := storage.NewPlayerPostgresRepository(tx)
		sessionRepository := storage.NewSessionPostgresRepository(tx)
		balanceRepository := storage.NewBalancePostgresRepository(tx)

		if _, err := playerRepository.Create(ctx, storage.CreatePlayerInput{
			ID:          playerID,
			DisplayName: displayName,
		}); err != nil {
			return err
		}

		session, err := sessionRepository.Create(ctx, storage.CreateSessionInput{
			ID:               sessionID,
			PlayerID:         playerID,
			SessionTokenHash: sessionTokenHash,
			ExpiresAt:        expiresAt,
		})
		if err != nil {
			return err
		}

		if _, err := balanceRepository.Create(ctx, storage.CreateBalanceInput{
			PlayerID:         playerID,
			AvailableBalance: s.config.Balance.StartingBalance,
			LockedBalance:    "0",
		}); err != nil {
			return err
		}

		createdSession = session
		return nil
	}); err != nil {
		return IssuedSession{}, fmt.Errorf("create guest session transaction: %w", err)
	}

	return IssuedSession{
		Session: Session{
			ID:               createdSession.ID,
			PlayerID:         createdSession.PlayerID,
			SessionTokenHash: createdSession.SessionTokenHash,
			ExpiresAt:        createdSession.ExpiresAt,
			CreatedAt:        createdSession.CreatedAt,
		},
		RawToken: rawToken,
	}, nil
}

func (s *service) ValidateSession(ctx context.Context, rawToken string) (AuthContext, error) {
	if rawToken == "" {
		return AuthContext{}, ErrUnauthorized
	}

	session, err := s.sessionRepository.GetByTokenHash(ctx, HashToken(rawToken))
	if err != nil {
		return AuthContext{}, fmt.Errorf("lookup session: %w", err)
	}

	if !session.ExpiresAt.After(domain.NowUTC()) {
		return AuthContext{}, ErrUnauthorized
	}

	return AuthContext{
		PlayerID:  session.PlayerID,
		SessionID: session.ID,
	}, nil
}
