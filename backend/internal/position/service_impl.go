package position

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"prediction/internal/balance"
	"prediction/internal/domain"
	"prediction/internal/market"
	"prediction/internal/storage"

	"github.com/jackc/pgx/v5"
)

type service struct {
	positionRepository storage.PositionRepository
	txManager          *storage.TxManager
	validator          Validator
}

type ServiceDeps struct {
	PositionRepository storage.PositionRepository
	TxManager          *storage.TxManager
	Validator          Validator
}

func NewService(deps ServiceDeps) Service {
	return &service{
		positionRepository: deps.PositionRepository,
		txManager:          deps.TxManager,
		validator:          deps.Validator,
	}
}

func (s *service) Create(ctx context.Context, playerID domain.PlayerID, input CreateInput) (Record, error) {
	normalized := normalizeCreateInput(input)
	if err := s.validator.ValidateCreateInput(ctx, playerID, normalized); err != nil {
		return Record{}, err
	}

	positionID, err := NewPositionID()
	if err != nil {
		return Record{}, err
	}

	potentialPayout, err := calculatePotentialPayout(normalized.StakeAmount)
	if err != nil {
		return Record{}, err
	}

	var created storage.Position
	if err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		marketRepository := storage.NewMarketPostgresRepository(tx)
		balanceRepository := storage.NewBalancePostgresRepository(tx)
		positionRepository := storage.NewPositionPostgresRepository(tx)

		selectedMarket, err := marketRepository.GetByID(ctx, normalized.MarketID)
		if err != nil {
			return fmt.Errorf("get market for position placement: %w", err)
		}

		if err := validateMarketPlacementState(selectedMarket); err != nil {
			return err
		}

		if _, err := balanceRepository.LockStake(ctx, storage.LockStakeInput{
			PlayerID: playerID,
			Amount:   normalized.StakeAmount,
		}); err != nil {
			return fmt.Errorf("lock stake for position placement: %w", err)
		}

		created, err = positionRepository.Create(ctx, storage.CreatePositionInput{
			ID:              positionID,
			PlayerID:        playerID,
			MarketID:        normalized.MarketID,
			Side:            normalized.Side,
			StakeAmount:     normalized.StakeAmount,
			PotentialPayout: potentialPayout,
		})
		if err != nil {
			return fmt.Errorf("create position: %w", err)
		}

		return nil
	}); err != nil {
		return Record{}, fmt.Errorf("place position: %w", err)
	}

	return toRecord(created), nil
}

func (s *service) ListByPlayerID(ctx context.Context, playerID domain.PlayerID, filter ListFilter) ([]Record, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	items, err := s.positionRepository.ListByPlayerID(ctx, playerID, filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("list positions by player id: %w", err)
	}

	records := make([]Record, 0, len(items))
	for _, item := range items {
		records = append(records, toRecord(item))
	}

	return records, nil
}

func (s *service) ValidateCreateInput(ctx context.Context, playerID domain.PlayerID, input CreateInput) error {
	return s.validator.ValidateCreateInput(ctx, playerID, normalizeCreateInput(input))
}

func normalizeCreateInput(input CreateInput) CreateInput {
	input.MarketID = domain.MarketID(strings.TrimSpace(string(input.MarketID)))
	input.StakeAmount = strings.TrimSpace(input.StakeAmount)
	return input
}

func toRecord(item storage.Position) Record {
	return Record{
		ID:              item.ID,
		PlayerID:        item.PlayerID,
		MarketID:        item.MarketID,
		Side:            item.Side,
		StakeAmount:     item.StakeAmount,
		PotentialPayout: item.PotentialPayout,
		Status:          item.Status,
		CreatedAt:       item.CreatedAt,
		SettledAt:       item.SettledAt,
	}
}

type ValidationDeps struct {
	MarketController  market.Controller
	BalanceController balance.Controller
}

func NewValidationService(deps ValidationDeps) Validator {
	return &validationService{
		marketController:  deps.MarketController,
		balanceController: deps.BalanceController,
	}
}

func validateMarketPlacementState(selectedMarket storage.Market) error {
	if selectedMarket.Status != domain.MarketStatusActive {
		return domain.NewValidationError("market_id", "market is not open for new positions", selectedMarket.ID)
	}

	if !selectedMarket.ExpiryTime.After(domain.NowUTC()) {
		return domain.NewValidationError("market_id", "market has already expired", selectedMarket.ID)
	}

	return nil
}

func calculatePotentialPayout(stakeAmount string) (string, error) {
	stake, err := parseDecimal(stakeAmount)
	if err != nil {
		return "", domain.NewValidationError("stake_amount", "stake amount must be a valid decimal value", stakeAmount)
	}

	payout := new(big.Rat).Mul(stake, big.NewRat(2, 1))
	return formatFixedScaleDecimal(payout, 8), nil
}
