package enrollment

import (
	"encoding/json"
	"errors"
	"net/http"
)

type asyncBody struct {
	CourseID int64 `json:"courseId"`
	UserID   int64 `json:"userId"`
}

// Handler exposes async enroll HTTP endpoints.
type Handler struct {
	enroller *Enroller
}

// NewHandler wraps an Enroller.
func NewHandler(e *Enroller) Handler {
	return Handler{enroller: e}
}

// CreateAsync handles POST /api/enrollments — validate + enqueue → 202 / 409 / 503.
func (h Handler) CreateAsync(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in asyncBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil || in.CourseID <= 0 || in.UserID <= 0 {
		http.Error(w, `{"code":"VALIDATION","message":"courseId and userId required"}`, http.StatusBadRequest)
		return
	}
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		http.Error(w, `{"code":"VALIDATION","message":"missing Idempotency-Key"}`, http.StatusBadRequest)
		return
	}
	switch err := h.enroller.Enqueue(Job{
		IdempotencyKey: key,
		CourseID:       in.CourseID,
		UserID:         in.UserID,
	}); {
	case err == nil:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"status":"accepted"}`))
	case errors.Is(err, ErrDuplicate):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"code":"DUPLICATE","message":"already accepted"}`))
	default: // ErrQueueFull or ErrShuttingDown
		w.Header().Set("Retry-After", "2")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"code":"QUEUE_FULL","message":"busy, retry later"}`))
	}
}
