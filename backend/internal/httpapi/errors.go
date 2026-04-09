package httpapi

import (
	"errors"
	"log"
	"net/http"

	"prediction/internal/auth"
	"prediction/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not_found"})
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		log.Printf("http internal error method=%s path=%s error=%v", r.Method, r.URL.Path, err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
	}
}
