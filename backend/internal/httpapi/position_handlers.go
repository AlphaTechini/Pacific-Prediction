package httpapi

import (
	"net/http"

	"prediction/internal/auth"
	"prediction/internal/domain"
	"prediction/internal/position"
)

type createPositionRequest struct {
	Side        string `json:"side"`
	StakeAmount string `json:"stake_amount"`
}

type listPositionsResponse struct {
	Positions []positionResponse `json:"positions"`
}

type positionResponse struct {
	ID              string  `json:"id"`
	PlayerID        string  `json:"player_id"`
	MarketID        string  `json:"market_id"`
	Side            string  `json:"side"`
	StakeAmount     string  `json:"stake_amount"`
	PotentialPayout string  `json:"potential_payout"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	SettledAt       *string `json:"settled_at,omitempty"`
}

func NewCreatePositionHandler(controller position.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID, err := auth.RequiredPlayerID(r.Context())
		if err != nil {
			writeError(w, r, err)
			return
		}

		marketID := domain.MarketID(r.PathValue("market_id"))
		if marketID == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_market_id"})
			return
		}

		var request createPositionRequest
		if err := decodeJSON(r, &request); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
			return
		}

		record, err := controller.Create(r.Context(), playerID, position.CreateInput{
			MarketID:    marketID,
			Side:        domain.PositionSide(request.Side),
			StakeAmount: request.StakeAmount,
		})
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusCreated, toPositionResponse(record))
	})
}

func NewListPlayerPositionsHandler(controller position.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID, err := auth.RequiredPlayerID(r.Context())
		if err != nil {
			writeError(w, r, err)
			return
		}

		records, err := controller.ListByPlayerID(r.Context(), playerID, position.ListFilter{
			Limit: 100,
		})
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusOK, listPositionsResponse{
			Positions: toPositionResponses(records),
		})
	})
}

func toPositionResponses(records []position.Record) []positionResponse {
	items := make([]positionResponse, 0, len(records))
	for _, record := range records {
		items = append(items, toPositionResponse(record))
	}

	return items
}

func toPositionResponse(record position.Record) positionResponse {
	response := positionResponse{
		ID:              string(record.ID),
		PlayerID:        string(record.PlayerID),
		MarketID:        string(record.MarketID),
		Side:            string(record.Side),
		StakeAmount:     record.StakeAmount,
		PotentialPayout: record.PotentialPayout,
		Status:          string(record.Status),
		CreatedAt:       record.CreatedAt.Format(http.TimeFormat),
	}

	if record.SettledAt != nil {
		value := record.SettledAt.Format(http.TimeFormat)
		response.SettledAt = &value
	}

	return response
}
