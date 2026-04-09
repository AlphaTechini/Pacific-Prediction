package httpapi

import (
	"net/http"

	"prediction/internal/auth"
	"prediction/internal/player"
)

type getMeResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewGetMeHandler(controller player.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID, err := auth.RequiredPlayerID(r.Context())
		if err != nil {
			writeError(w, r, err)
			return
		}

		profile, err := controller.GetMe(r.Context(), playerID)
		if err != nil {
			writeError(w, r, err)
			return
		}

		writeJSON(w, http.StatusOK, getMeResponse{
			ID:          string(profile.ID),
			DisplayName: profile.DisplayName,
			CreatedAt:   profile.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:   profile.UpdatedAt.Format(http.TimeFormat),
		})
	})
}
