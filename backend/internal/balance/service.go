package balance

import (
	"context"
	"fmt"

	"prediction/internal/domain"
	"prediction/internal/storage"
)

type service struct {
	balanceRepository storage.BalanceRepository
}

func NewService(balanceRepository storage.BalanceRepository) Service {
	return &service{balanceRepository: balanceRepository}
}

func (s *service) GetBalance(ctx context.Context, playerID domain.PlayerID) (Snapshot, error) {
	balance, err := s.balanceRepository.GetByPlayerID(ctx, playerID)
	if err != nil {
		return Snapshot{}, fmt.Errorf("get balance: %w", err)
	}

	return Snapshot{
		PlayerID:         balance.PlayerID,
		AvailableBalance: balance.AvailableBalance,
		LockedBalance:    balance.LockedBalance,
		UpdatedAt:        balance.UpdatedAt,
	}, nil
}

func (s *service) LockStake(ctx context.Context, playerID domain.PlayerID, amount string) error {
	if _, err := s.balanceRepository.LockStake(ctx, storage.LockStakeInput{
		PlayerID: playerID,
		Amount:   amount,
	}); err != nil {
		return fmt.Errorf("lock stake: %w", err)
	}

	return nil
}

func (s *service) UnlockStake(ctx context.Context, playerID domain.PlayerID, amount string) error {
	return domain.NewValidationError("amount", "unlock stake is not implemented yet", amount)
}

func (s *service) Credit(ctx context.Context, playerID domain.PlayerID, amount string) error {
	return domain.NewValidationError("amount", "credit is not implemented yet", amount)
}

func (s *service) Debit(ctx context.Context, playerID domain.PlayerID, amount string) error {
	return domain.NewValidationError("amount", "debit is not implemented yet", amount)
}
