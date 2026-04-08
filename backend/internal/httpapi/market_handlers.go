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
	Title             string `json:"title"`
	Symbol            string `json:"symbol"`
	MarketType        string `json:"market_type"`
	ConditionOperator string `json:"condition_operator"`
	ThresholdValue    string `json:"threshold_value"`
	SourceType        string `json:"source_type"`
	SourceInterval    string `json:"source_interval"`
	ReferenceValue    string `json:"reference_value"`
	ExpiryTime        string `json:"expiry_time"`
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
			writeError(w, err)
			return
		}

		var request createMarketRequest
		if err := decodeJSON(r, &request); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
			return
		}

		expiryTime, err := time.Parse(time.RFC3339, request.ExpiryTime)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "expiry_time must be RFC3339"})
			return
		}

		record, err := controller.Create(r.Context(), market.CreateInput{
			Title:             request.Title,
			Symbol:            request.Symbol,
			MarketType:        domain.MarketType(request.MarketType),
			ConditionOperator: domain.ConditionOperator(request.ConditionOperator),
			ThresholdValue:    request.ThresholdValue,
			SourceType:        domain.SourceType(request.SourceType),
			SourceInterval:    request.SourceInterval,
			ReferenceValue:    request.ReferenceValue,
			ExpiryTime:        expiryTime,
			CreatedByPlayerID: playerID,
		})
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, toCreateMarketResponse(record))
	})
}

func NewListMarketsHandler(controller market.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		catalog, err := controller.ListCatalog(r.Context(), 50)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, listMarketsResponse{
			Active:   toMarketResponses(catalog.Active),
			Resolved: toMarketResponses(catalog.Resolved),
		})
	})
}

func NewGetMarketDetailHandler(controller market.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "/api/v1/markets/"

		marketID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
		if marketID == "" || strings.Contains(marketID, "/") {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_market_id"})
			return
		}

		record, err := controller.GetByID(r.Context(), domain.MarketID(marketID))
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, toMarketResponse(record))
	})
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
