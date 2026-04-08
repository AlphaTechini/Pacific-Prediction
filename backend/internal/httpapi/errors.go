package httpapi

import (
	"errors"
	"net/http"

	"prediction/internal/auth"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
	}
}
