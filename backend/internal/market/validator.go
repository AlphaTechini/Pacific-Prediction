package market

import (
	"context"
	"fmt"
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

type ValidationConfig struct {
	PriceThresholdCreationBandPercent string
}

type validationService struct {
	symbolProvider SymbolProvider
	config         ValidationConfig
	now            func() time.Time
}

func NewValidationService(symbolProvider SymbolProvider, config ValidationConfig) Validator {
	return &validationService{
		symbolProvider: symbolProvider,
		config:         config,
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

	if input.CreatorSide != domain.PositionSideYes && input.CreatorSide != domain.PositionSideNo {
		return domain.NewValidationError("creator_side", "creator side must be yes or no", input.CreatorSide)
	}

	creatorStake, err := domain.ParseDecimal(input.CreatorStakeAmount)
	if err != nil {
		return domain.NewValidationError("creator_stake_amount", "creator stake amount must be a valid decimal value", input.CreatorStakeAmount)
	}
	if creatorStake.Sign() <= 0 {
		return domain.NewValidationError("creator_stake_amount", "creator stake amount must be greater than zero", input.CreatorStakeAmount)
	}
	if !domain.FitsNumericScale(input.CreatorStakeAmount, 8) {
		return domain.NewValidationError("creator_stake_amount", "creator stake amount must use no more than 8 decimal places", input.CreatorStakeAmount)
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

	switch input.MarketType {
	case domain.MarketTypePriceThreshold:
		return validatePriceThresholdInput(s.now(), input, s.config)
	case domain.MarketTypeCandleDirection:
		return validateCandleDirectionInput(s.now(), input, model)
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

func validatePriceThresholdInput(now time.Time, input CreateInput, config ValidationConfig) error {
	if err := validateExpiryTime(now, input.ExpiryTime); err != nil {
		return err
	}

	if input.SourceInterval != "" {
		return domain.NewValidationError("source_interval", "price threshold markets do not use source intervals", input.SourceInterval)
	}

	if err := requireDecimal("threshold_value", input.ThresholdValue); err != nil {
		return err
	}

	if !domain.FitsNumericScale(input.ThresholdValue, 8) {
		return domain.NewValidationError("threshold_value", "threshold value must use no more than 8 decimal places", input.ThresholdValue)
	}

	if err := validateThresholdScaleForSymbol(input); err != nil {
		return err
	}

	return validatePriceThresholdRange(input, config)
}

func validateCandleDirectionInput(now time.Time, input CreateInput, model ValidationModel) error {
	if input.ThresholdValue != "" {
		return domain.NewValidationError("threshold_value", "candle direction markets do not use threshold values", input.ThresholdValue)
	}

	if strings.TrimSpace(input.SourceInterval) == "" {
		return domain.NewValidationError("source_interval", "candle direction markets require a source interval", input.SourceInterval)
	}

	if !containsString(model.AllowedIntervals, input.SourceInterval) {
		return domain.NewValidationError("source_interval", "unsupported candle interval", input.SourceInterval)
	}

	if err := validateExpiryTime(now, input.ExpiryTime); err != nil {
		return err
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

func validateThresholdScaleForSymbol(input CreateInput) error {
	maxScale := 8
	if strings.TrimSpace(input.SymbolTickSize) != "" {
		tickScale := domain.DecimalScale(input.SymbolTickSize)
		if tickScale < maxScale {
			maxScale = tickScale
		}
	}

	if !domain.FitsNumericScale(input.ThresholdValue, maxScale) {
		return domain.NewValidationError(
			"threshold_value",
			fmt.Sprintf("threshold value must use no more than %d decimal places for this symbol", maxScale),
			input.ThresholdValue,
		)
	}

	return nil
}

func validatePriceThresholdRange(input CreateInput, config ValidationConfig) error {
	if strings.TrimSpace(config.PriceThresholdCreationBandPercent) == "" {
		return fmt.Errorf("price-threshold creation band config is required")
	}

	threshold, err := domain.ParseDecimal(input.ThresholdValue)
	if err != nil {
		return domain.NewValidationError("threshold_value", "threshold value must be a valid decimal string", input.ThresholdValue)
	}

	referenceValue, err := domain.ParseDecimal(input.ReferenceValue)
	if err != nil {
		return domain.NewValidationError("reference_value", "reference value must be a valid decimal string", input.ReferenceValue)
	}

	if referenceValue.Sign() <= 0 {
		return domain.NewValidationError("reference_value", "reference value must be greater than zero", input.ReferenceValue)
	}

	tickSize, err := domain.ParseDecimal(input.SymbolTickSize)
	if err != nil {
		return domain.NewValidationError("symbol", "symbol tick size must be a valid decimal string", input.SymbolTickSize)
	}

	if tickSize.Sign() <= 0 {
		return domain.NewValidationError("symbol", "symbol tick size must be greater than zero", input.SymbolTickSize)
	}

	bandPercent, err := domain.ParseDecimal(config.PriceThresholdCreationBandPercent)
	if err != nil {
		return fmt.Errorf("parse price-threshold creation band percent: %w", err)
	}

	lowerBound, upperBound := calculatePriceThresholdBounds(referenceValue, bandPercent)
	displayScale := thresholdDisplayScale(input.SymbolTickSize)
	referenceDisplay := domain.FormatFixedScaleDecimal(referenceValue, displayScale)

	distance := new(big.Rat).Sub(new(big.Rat).Set(threshold), referenceValue)
	if distance.Sign() < 0 {
		distance.Neg(distance)
	}
	if distance.Cmp(tickSize) < 0 {
		return domain.NewValidationError(
			"threshold_value",
			fmt.Sprintf("threshold value must be at least one tick away from the creation reference %s", referenceDisplay),
			input.ThresholdValue,
		)
	}

	switch input.ConditionOperator {
	case domain.ConditionOperatorGT, domain.ConditionOperatorGTE:
		if threshold.Cmp(referenceValue) <= 0 {
			return domain.NewValidationError(
				"threshold_value",
				fmt.Sprintf("threshold value must be above the creation reference %s", referenceDisplay),
				input.ThresholdValue,
			)
		}

		if threshold.Cmp(upperBound) > 0 {
			return domain.NewValidationError(
				"threshold_value",
				fmt.Sprintf(
					"threshold value must stay within the %s%% upper band from the creation reference (%s)",
					config.PriceThresholdCreationBandPercent,
					domain.FormatFixedScaleDecimal(upperBound, displayScale),
				),
				input.ThresholdValue,
			)
		}
	case domain.ConditionOperatorLT, domain.ConditionOperatorLTE:
		if threshold.Cmp(referenceValue) >= 0 {
			return domain.NewValidationError(
				"threshold_value",
				fmt.Sprintf("threshold value must be below the creation reference %s", referenceDisplay),
				input.ThresholdValue,
			)
		}

		if threshold.Cmp(lowerBound) < 0 {
			return domain.NewValidationError(
				"threshold_value",
				fmt.Sprintf(
					"threshold value must stay within the %s%% lower band from the creation reference (%s)",
					config.PriceThresholdCreationBandPercent,
					domain.FormatFixedScaleDecimal(lowerBound, displayScale),
				),
				input.ThresholdValue,
			)
		}
	default:
		return domain.NewValidationError("condition_operator", "operator is not supported for price-threshold range validation", input.ConditionOperator)
	}

	return nil
}

func calculatePriceThresholdBounds(referenceValue, bandPercent *big.Rat) (*big.Rat, *big.Rat) {
	oneHundred := big.NewRat(100, 1)
	bandRatio := new(big.Rat).Quo(new(big.Rat).Set(bandPercent), oneHundred)
	lowerMultiplier := new(big.Rat).Sub(big.NewRat(1, 1), bandRatio)
	upperMultiplier := new(big.Rat).Add(big.NewRat(1, 1), bandRatio)

	lowerBound := new(big.Rat).Mul(new(big.Rat).Set(referenceValue), lowerMultiplier)
	upperBound := new(big.Rat).Mul(new(big.Rat).Set(referenceValue), upperMultiplier)

	return lowerBound, upperBound
}

func thresholdDisplayScale(tickSize string) int {
	scale := domain.DecimalScale(tickSize)
	if scale <= 0 {
		return 0
	}
	if scale > 8 {
		return 8
	}

	return scale
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
