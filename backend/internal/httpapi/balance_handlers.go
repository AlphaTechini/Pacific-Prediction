package httpapi

import (
	"net/http"

	"prediction/internal/auth"
	"prediction/internal/balance"
)

type getBalanceResponse struct {
	PlayerID         string `json:"player_id"`
	AvailableBalance string `json:"available_balance"`
	LockedBalance    string `json:"locked_balance"`
	UpdatedAt        string `json:"updated_at"`
}

func NewGetBalanceHandler(controller balance.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID, err := auth.RequiredPlayerID(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}

		snapshot, err := controller.GetBalance(r.Context(), playerID)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, getBalanceResponse{
			PlayerID:         string(snapshot.PlayerID),
			AvailableBalance: snapshot.AvailableBalance,
			LockedBalance:    snapshot.LockedBalance,
			UpdatedAt:        snapshot.UpdatedAt.Format(http.TimeFormat),
		})
	})
}
