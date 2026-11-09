package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

// Middleware wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

type ctxKey string

const requestIDKey ctxKey = "requestID"

// RequestIDFrom returns the request id attached by Logging middleware.
func RequestIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

// Logging emits structured JSON access logs with request_id.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rid := r.Header.Get("X-Request-ID")
		if !validRequestID(rid) {
			rid = newRequestID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", rid)

		sw := &statusWriter{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(sw, r)
		slog.Info("http_request",
			"request_id", rid,
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.code,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// Recovery catches panics in the handler goroutine and returns 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic", "request_id", RequestIDFrom(r.Context()), "err", rec)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// validRequestID accepts a short id of letters, digits, '-' and '_'.
// Anything else (empty, too long, spaces, control characters) is treated as missing.
func validRequestID(s string) bool {
	if s == "" || len(s) > 64 {
		return false
	}
	for _, c := range s {
		ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_'
		if !ok {
			return false
		}
	}
	return true
}

func newRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
