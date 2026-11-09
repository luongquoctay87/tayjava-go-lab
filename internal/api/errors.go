package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/tayjava/go-lab/internal/course"
)

// ErrorBody is the stable JSON error contract for API clients.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSON writes a JSON response with the given status.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes a structured ErrorBody.
func WriteError(w http.ResponseWriter, status int, code, msg string) {
	WriteJSON(w, status, ErrorBody{Code: code, Message: msg})
}

// MapDomainError maps domain errors to HTTP status + ErrorBody.
// Internal details are logged server-side only.
func MapDomainError(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		return
	case errors.Is(err, course.ErrNotFound):
		WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, course.ErrNoSeats):
		WriteError(w, http.StatusConflict, "NO_SEATS", err.Error())
	case errors.Is(err, course.ErrConflict):
		WriteError(w, http.StatusConflict, "CONFLICT", err.Error())
	default:
		// Log the original error server-side; do not leak internal details to the client.
		slog.Error("unhandled error", "err", err)
		WriteError(w, http.StatusInternalServerError, "INTERNAL", "internal error")
	}
}
