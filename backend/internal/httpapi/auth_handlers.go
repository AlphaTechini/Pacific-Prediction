package httpapi

import (
	"errors"
	"net/http"

	"prediction/internal/auth"
	"prediction/internal/player"
)

type createGuestSessionRequest struct {
	DisplayName string `json:"display_name"`
}

type createGuestSessionResponse struct {
	PlayerID    string `json:"player_id"`
	DisplayName string `json:"display_name"`
	ExpiresAt   string `json:"expires_at"`
}

func NewCreateGuestSessionHandler(authController auth.Controller, playerController player.Controller, cookies *auth.CookieManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request createGuestSessionRequest
		if r.ContentLength > 0 {
			if err := decodeJSON(r, &request); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
				return
			}
		}

		issuedSession, err := authController.CreateGuestSession(r.Context(), auth.CreateGuestSessionInput{
			DisplayName: request.DisplayName,
		})
		if err != nil {
			writeError(w, r, err)
			return
		}

		profile, err := playerController.GetMe(r.Context(), issuedSession.Session.PlayerID)
		if err != nil {
			writeError(w, r, err)
			return
		}

		cookies.SetSessionCookie(w, issuedSession.RawToken, issuedSession.Session.ExpiresAt)

		writeJSON(w, http.StatusCreated, createGuestSessionResponse{
			PlayerID:    string(profile.ID),
			DisplayName: profile.DisplayName,
			ExpiresAt:   issuedSession.Session.ExpiresAt.Format(http.TimeFormat),
		})
	})
}

func NewRequireSessionMiddleware(controller auth.Controller, cookies *auth.CookieManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken, err := cookies.ReadSessionCookie(r)
			if err != nil {
				writeError(w, r, auth.ErrUnauthorized)
				return
			}

			authContext, err := controller.ValidateSession(r.Context(), rawToken)
			if err != nil {
				if errors.Is(err, auth.ErrUnauthorized) {
					writeError(w, r, auth.ErrUnauthorized)
					return
				}

				writeError(w, r, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithAuthContext(r.Context(), authContext)))
		})
	}
}
