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

// CreateAsync handles POST /api/enrollments — validate + enqueue → 202 / 503.
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
	err := h.enroller.Enqueue(r.Context(), Job{
		IdempotencyKey: key,
		CourseID:       in.CourseID,
		UserID:         in.UserID,
	})
	if errors.Is(err, ErrQueueFull) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"code":"QUEUE_FULL","message":"enrollment queue full"}`))
		return
	}
	if err != nil {
		http.Error(w, "canceled", http.StatusRequestTimeout)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"status":"accepted"}`))
}
