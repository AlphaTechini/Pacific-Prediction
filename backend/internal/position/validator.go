package position

import (
	"context"
	"fmt"
	"strings"

	"prediction/internal/balance"
	"prediction/internal/domain"
	"prediction/internal/market"
)

type Validator interface {
	ValidateCreateInput(ctx context.Context, playerID domain.PlayerID, input CreateInput) error
}

type validationService struct {
	marketController  market.Controller
	balanceController balance.Controller
}

func (s *validationService) ValidateCreateInput(ctx context.Context, playerID domain.PlayerID, input CreateInput) error {
	if input.MarketID == "" {
		return domain.NewValidationError("market_id", "market id is required", input.MarketID)
	}

	if input.Side != domain.PositionSideYes && input.Side != domain.PositionSideNo {
		return domain.NewValidationError("side", "side must be yes or no", input.Side)
	}

	stake, err := domain.ParseDecimal(input.StakeAmount)
	if err != nil {
		return domain.NewValidationError("stake_amount", "stake amount must be a valid decimal value", input.StakeAmount)
	}
	if stake.Sign() <= 0 {
		return domain.NewValidationError("stake_amount", "stake amount must be greater than zero", input.StakeAmount)
	}
	if !domain.IsWholeNumber(strings.TrimSpace(input.StakeAmount)) {
		return domain.NewValidationError("stake_amount", "stake amount must be a whole number", input.StakeAmount)
	}

	selectedMarket, err := s.marketController.GetByID(ctx, input.MarketID)
	if err != nil {
		return fmt.Errorf("get market for position validation: %w", err)
	}

	if selectedMarket.Status != domain.MarketStatusActive {
		return domain.NewValidationError("market_id", "market is not open for new positions", input.MarketID)
	}

	if !selectedMarket.ExpiryTime.After(domain.NowUTC()) {
		return domain.NewValidationError("market_id", "market has already expired", input.MarketID)
	}

	balanceSnapshot, err := s.balanceController.GetBalance(ctx, playerID)
	if err != nil {
		return fmt.Errorf("get balance for position validation: %w", err)
	}

	availableBalance, err := domain.ParseDecimal(balanceSnapshot.AvailableBalance)
	if err != nil {
		return fmt.Errorf("parse available balance: %w", err)
	}

	if stake.Cmp(availableBalance) > 0 {
		return domain.NewValidationError("stake_amount", "insufficient available balance", input.StakeAmount)
	}

	return nil
}
