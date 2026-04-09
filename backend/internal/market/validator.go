package market

import (
	"context"
	"math/big"
	"strings"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
)

type SymbolProvider interface {
	ListMarketInfo(ctx context.Context) ([]pacifica.MarketInfo, error)
}

type Validator interface {
	ValidateCreateInput(ctx context.Context, input CreateInput) error
}

type validationService struct {
	symbolProvider SymbolProvider
	now            func() time.Time
}

func NewValidationService(symbolProvider SymbolProvider) Validator {
	return &validationService{
		symbolProvider: symbolProvider,
		now: func() time.Time {
			return domain.NowUTC()
		},
	}
}

func (s *validationService) ValidateCreateInput(ctx context.Context, input CreateInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return domain.NewValidationError("title", "title is required", input.Title)
	}

	if input.CreatedByPlayerID == "" {
		return domain.NewValidationError("created_by_player_id", "creator is required", input.CreatedByPlayerID)
	}

	model, ok := findValidationModel(input.MarketType)
	if !ok {
		return domain.NewValidationError("market_type", "unsupported market type", input.MarketType)
	}

	if strings.TrimSpace(input.Symbol) == "" {
		return domain.NewValidationError("symbol", "symbol is required", input.Symbol)
	}

	if err := s.validateSymbol(ctx, input.Symbol); err != nil {
		return err
	}

	if input.SourceType != model.SourceType {
		return domain.NewValidationError("source_type", "source type does not match market type", input.SourceType)
	}

	if !containsOperator(model.AllowedOperators, input.ConditionOperator) {
		return domain.NewValidationError("condition_operator", "operator is not supported for market type", input.ConditionOperator)
	}

	if err := validateExpiryTime(s.now(), input.ExpiryTime); err != nil {
		return err
	}

	switch input.MarketType {
	case domain.MarketTypePriceThreshold:
		return validatePriceThresholdInput(input)
	case domain.MarketTypeCandleDirection:
		return validateCandleDirectionInput(input, model)
	case domain.MarketTypeFundingThreshold:
		return validateFundingInput(input, model)
	default:
		return domain.NewValidationError("market_type", "unsupported market type", input.MarketType)
	}
}

func (s *validationService) validateSymbol(ctx context.Context, symbol string) error {
	items, err := s.symbolProvider.ListMarketInfo(ctx)
	if err != nil {
		return err
	}

	target := strings.ToUpper(strings.TrimSpace(symbol))
	for _, item := range items {
		if strings.EqualFold(item.Symbol, target) {
			return nil
		}
	}

	return domain.NewValidationError("symbol", "symbol is not supported by Pacifica", symbol)
}

func findValidationModel(marketType domain.MarketType) (ValidationModel, bool) {
	for _, model := range SupportedValidationModels() {
		if model.MarketType == marketType {
			return model, true
		}
	}

	return ValidationModel{}, false
}

func validateExpiryTime(now, expiry time.Time) error {
	if expiry.IsZero() {
		return domain.NewValidationError("expiry_time", "expiry time is required", expiry)
	}

	if !expiry.After(now) {
		return domain.NewValidationError("expiry_time", "expiry time must be in the future", expiry)
	}

	return nil
}

func validatePriceThresholdInput(input CreateInput) error {
	if input.SourceInterval != "" {
		return domain.NewValidationError("source_interval", "price threshold markets do not use source intervals", input.SourceInterval)
	}

	if err := requireDecimal("threshold_value", input.ThresholdValue); err != nil {
		return err
	}

	return nil
}

func validateCandleDirectionInput(input CreateInput, model ValidationModel) error {
	if input.ThresholdValue != "" {
		return domain.NewValidationError("threshold_value", "candle direction markets do not use threshold values", input.ThresholdValue)
	}

	if strings.TrimSpace(input.SourceInterval) == "" {
		return domain.NewValidationError("source_interval", "candle direction markets require a source interval", input.SourceInterval)
	}

	if !containsString(model.AllowedIntervals, input.SourceInterval) {
		return domain.NewValidationError("source_interval", "unsupported candle interval", input.SourceInterval)
	}

	if _, _, err := domain.CandleWindowForExpiry(input.ExpiryTime, input.SourceInterval); err != nil {
		return domain.NewValidationError("expiry_time", "candle market expiry must align to the selected candle close boundary", input.ExpiryTime)
	}

	return nil
}

func validateFundingInput(input CreateInput, model ValidationModel) error {
	if strings.TrimSpace(input.SourceInterval) != domain.SourceIntervalFundingEpoch {
		return domain.NewValidationError("source_interval", "funding markets require funding_epoch interval", input.SourceInterval)
	}

	switch input.ConditionOperator {
	case domain.ConditionOperatorPositive, domain.ConditionOperatorNegative:
		if strings.TrimSpace(input.ThresholdValue) != "" {
			return domain.NewValidationError("threshold_value", "sign-based funding markets should not include a threshold", input.ThresholdValue)
		}
	default:
		if err := requireDecimal("threshold_value", input.ThresholdValue); err != nil {
			return err
		}
	}

	if model.SourceType != input.SourceType {
		return domain.NewValidationError("source_type", "source type does not match funding market rules", input.SourceType)
	}

	return nil
}

func requireDecimal(field, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return domain.NewValidationError(field, "value is required", value)
	}

	number, ok := new(big.Rat).SetString(trimmed)
	if !ok {
		return domain.NewValidationError(field, "value must be a valid decimal string", value)
	}

	if number.Sign() == 0 {
		return domain.NewValidationError(field, "value must be non-zero", value)
	}

	return nil
}

func containsOperator(items []domain.ConditionOperator, target domain.ConditionOperator) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}

	return false
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}

	return false
}
