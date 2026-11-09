package api

import (
	"bytes"
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
	if err == nil {
		return
	}
	status, body := domainErrorResponse(err)
	if status == http.StatusInternalServerError {
		slog.Error("internal_error", "err", err.Error())
	}
	WriteError(w, status, body.Code, body.Message)
}

func domainErrorResponse(err error) (int, ErrorBody) {
	switch {
	case errors.Is(err, course.ErrNotFound):
		return http.StatusNotFound, ErrorBody{Code: "NOT_FOUND", Message: err.Error()}
	case errors.Is(err, course.ErrNoSeats):
		return http.StatusConflict, ErrorBody{Code: "NO_SEATS", Message: err.Error()}
	case errors.Is(err, course.ErrConflict):
		return http.StatusConflict, ErrorBody{Code: "CONFLICT", Message: err.Error()}
	default:
		return http.StatusInternalServerError, ErrorBody{Code: "INTERNAL", Message: "internal error"}
	}
}

// jsonBody encodes v the same way WriteJSON does, including the trailing newline.
func jsonBody(v any) []byte {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(v)
	return buf.Bytes()
}
