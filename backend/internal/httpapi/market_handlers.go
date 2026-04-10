package httpapi

import (
	"net/http"
	"strings"
	"time"

	"prediction/internal/auth"
	"prediction/internal/domain"
	"prediction/internal/market"
)

type createMarketRequest struct {
	Title              string `json:"title"`
	Symbol             string `json:"symbol"`
	MarketType         string `json:"market_type"`
	ConditionOperator  string `json:"condition_operator"`
	CreatorSide        string `json:"creator_side"`
	CreatorStakeAmount string `json:"creator_stake_amount"`
	ThresholdValue     string `json:"threshold_value"`
	SourceType         string `json:"source_type"`
	SourceInterval     string `json:"source_interval"`
	ReferenceValue     string `json:"reference_value"`
	ExpiryTime         string `json:"expiry_time"`
}

type createMarketResponse struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Symbol            string  `json:"symbol"`
	MarketType        string  `json:"market_type"`
	ConditionOperator string  `json:"condition_operator"`
	ThresholdValue    string  `json:"threshold_value,omitempty"`
	SourceType        string  `json:"source_type"`
	SourceInterval    string  `json:"source_interval,omitempty"`
	ReferenceValue    string  `json:"reference_value,omitempty"`
	ExpiryTime        string  `json:"expiry_time"`
	Status            string  `json:"status"`
	Result            string  `json:"result,omitempty"`
	SettlementValue   string  `json:"settlement_value,omitempty"`
	ResolvedAt        *string `json:"resolved_at,omitempty"`
	ResolutionReason  string  `json:"resolution_reason,omitempty"`
	CreatedByPlayerID string  `json:"created_by_player_id"`
	CreatedAt         string  `json:"created_at"`
}

type listMarketsResponse struct {
	Active   []marketResponse `json:"active"`
	Resolved []marketResponse `json:"resolved"`
}

type marketCreateContextResponse struct {
	Symbols                           []marketCreateContextSymbolResponse `json:"symbols"`
	ValidationModels                  []marketValidationModelResponse     `json:"validation_models"`
	PriceThresholdCreationBandPercent string                              `json:"price_threshold_creation_band_percent"`
}

type marketCreateContextSymbolResponse struct {
	Symbol          string `json:"symbol"`
	TickSize        string `json:"tick_size"`
	MinTick         string `json:"min_tick"`
	MaxTick         string `json:"max_tick"`
	LotSize         string `json:"lot_size"`
	MinOrderSize    string `json:"min_order_size"`
	MaxOrderSize    string `json:"max_order_size"`
	MaxLeverage     int    `json:"max_leverage"`
	IsolatedOnly    bool   `json:"isolated_only"`
	MarkPrice       string `json:"mark_price,omitempty"`
	OraclePrice     string `json:"oracle_price,omitempty"`
	FundingRate     string `json:"funding_rate,omitempty"`
	NextFundingRate string `json:"next_funding_rate,omitempty"`
	OpenInterest    string `json:"open_interest,omitempty"`
	Volume24H       string `json:"volume_24h,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

type marketValidationModelResponse struct {
	MarketType        string   `json:"market_type"`
	SourceType        string   `json:"source_type"`
	AllowedOperators  []string `json:"allowed_operators"`
	RequiresThreshold bool     `json:"requires_threshold"`
	RequiresInterval  bool     `json:"requires_interval"`
	AllowedIntervals  []string `json:"allowed_intervals,omitempty"`
}

type marketResponse struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Symbol            string  `json:"symbol"`
	MarketType        string  `json:"market_type"`
	ConditionOperator string  `json:"condition_operator"`
	ThresholdValue    string  `json:"threshold_value,omitempty"`
	SourceType        string  `json:"source_type"`
	SourceInterval    string  `json:"source_interval,omitempty"`
	ReferenceValue    string  `json:"reference_value,omitempty"`
	ExpiryTime        string  `json:"expiry_time"`
	Status            string  `json:"status"`
	Result            string  `json:"result,omitempty"`
	SettlementValue   string  `json:"settlement_value,omitempty"`
	ResolvedAt        *string `json:"resolved_at,omitempty"`
	ResolutionReason  string  `json:"resolution_reason,omitempty"`
	CreatedByPlayerID string  `json:"created_by_player_id"`
	CreatedAt         string  `json:"created_at"`
}

func NewCreateMarketHandler(controller market.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID, err := auth.RequiredPlayerID(r.Context())
		if err != nil {
			writeError(w, r, err)
			return
		}

		var request createMarketRequest
		if err := decodeJSON(r, &request); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
			return
		}

		var expiryTime time.Time
		if strings.TrimSpace(request.ExpiryTime) != "" {
			expiryTime, err = time.Parse(time.RFC3339, request.ExpiryTime)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "expiry_time must be RFC3339"})
				return
			}
		}

		record, err := controller.Create(r.Context(), market.CreateInput{
			Title:              request.Title,
			Symbol:             request.Symbol,
			MarketType:         domain.MarketType(request.MarketType),
			ConditionOperator:  domain.ConditionOperator(request.ConditionOperator),
			CreatorSide:        domain.PositionSide(request.CreatorSide),
			CreatorStakeAmount: request.CreatorStakeAmount,
			ThresholdValue:     request.ThresholdValue,
			SourceType:         domain.SourceType(request.SourceType),
			SourceInterval:     request.SourceInterval,
			ReferenceValue:     request.ReferenceValue,
			ExpiryTime:         expiryTime,
			CreatedByPlayerID:  playerID,
		})
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusCreated, toCreateMarketResponse(record))
	})
}

func NewListMarketsHandler(controller market.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		catalog, err := controller.ListCatalog(r.Context(), 50)
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusOK, listMarketsResponse{
			Active:   toMarketResponses(catalog.Active),
			Resolved: toMarketResponses(catalog.Resolved),
		})
	})
}

func NewGetMarketCreateContextHandler(controller market.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createContext, err := controller.GetCreateContext(r.Context())
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusOK, toMarketCreateContextResponse(createContext))
	})
}

func NewGetMarketDetailHandler(controller market.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		marketID := r.PathValue("market_id")
		if marketID == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_market_id"})
			return
		}

		record, err := controller.GetByID(r.Context(), domain.MarketID(marketID))
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusOK, toMarketResponse(record))
	})
}

func toMarketCreateContextResponse(context market.CreateContext) marketCreateContextResponse {
	symbols := make([]marketCreateContextSymbolResponse, 0, len(context.Symbols))
	for _, item := range context.Symbols {
		response := marketCreateContextSymbolResponse{
			Symbol:          item.Symbol,
			TickSize:        item.TickSize,
			MinTick:         item.MinTick,
			MaxTick:         item.MaxTick,
			LotSize:         item.LotSize,
			MinOrderSize:    item.MinOrderSize,
			MaxOrderSize:    item.MaxOrderSize,
			MaxLeverage:     item.MaxLeverage,
			IsolatedOnly:    item.IsolatedOnly,
			MarkPrice:       item.MarkPrice,
			OraclePrice:     item.OraclePrice,
			FundingRate:     item.FundingRate,
			NextFundingRate: item.NextFundingRate,
			OpenInterest:    item.OpenInterest,
			Volume24H:       item.Volume24H,
		}
		if !item.UpdatedAt.IsZero() {
			response.UpdatedAt = item.UpdatedAt.Format(time.RFC3339)
		}
		symbols = append(symbols, response)
	}

	models := make([]marketValidationModelResponse, 0, len(context.ValidationModels))
	for _, item := range context.ValidationModels {
		operators := make([]string, 0, len(item.AllowedOperators))
		for _, operator := range item.AllowedOperators {
			operators = append(operators, string(operator))
		}

		models = append(models, marketValidationModelResponse{
			MarketType:        string(item.MarketType),
			SourceType:        string(item.SourceType),
			AllowedOperators:  operators,
			RequiresThreshold: item.RequiresThreshold,
			RequiresInterval:  item.RequiresInterval,
			AllowedIntervals:  item.AllowedIntervals,
		})
	}

	return marketCreateContextResponse{
		Symbols:                           symbols,
		ValidationModels:                  models,
		PriceThresholdCreationBandPercent: context.PriceThresholdCreationBandPercent,
	}
}

func toCreateMarketResponse(record market.Record) createMarketResponse {
	base := toMarketResponse(record)
	return createMarketResponse(base)
}

func toMarketResponses(records []market.Record) []marketResponse {
	items := make([]marketResponse, 0, len(records))
	for _, record := range records {
		items = append(items, toMarketResponse(record))
	}

	return items
}

func toMarketResponse(record market.Record) marketResponse {
	response := marketResponse{
		ID:                string(record.ID),
		Title:             record.Title,
		Symbol:            record.Symbol,
		MarketType:        string(record.MarketType),
		ConditionOperator: string(record.ConditionOperator),
		ThresholdValue:    record.ThresholdValue,
		SourceType:        string(record.SourceType),
		SourceInterval:    record.SourceInterval,
		ReferenceValue:    record.ReferenceValue,
		ExpiryTime:        record.ExpiryTime.Format(time.RFC3339),
		Status:            string(record.Status),
		Result:            string(record.Result),
		SettlementValue:   record.SettlementValue,
		ResolutionReason:  record.ResolutionReason,
		CreatedByPlayerID: string(record.CreatedByPlayerID),
		CreatedAt:         record.CreatedAt.Format(time.RFC3339),
	}

	if record.ResolvedAt != nil {
		value := record.ResolvedAt.Format(time.RFC3339)
		response.ResolvedAt = &value
	}

	return response
}
