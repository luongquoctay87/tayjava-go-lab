package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tayjava/go-lab/internal/course"
)

// Enroller is the part of the service the handler needs.
// *course.Service satisfies it, and tests can pass a smaller fake.
type Enroller interface {
	Enroll(ctx context.Context, courseID, userID int64) error
}

// EnrollmentHandler exposes enrollment HTTP endpoints.
type EnrollmentHandler struct {
	svc Enroller
}

// NewEnrollmentHandler constructs an EnrollmentHandler.
func NewEnrollmentHandler(svc Enroller) EnrollmentHandler {
	return EnrollmentHandler{svc: svc}
}

// Create handles POST /api/enrollments with a JSON body.
func (h EnrollmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in struct {
		CourseID int64 `json:"courseId"`
		UserID   int64 `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if in.CourseID == 0 || in.UserID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "courseId and userId required"})
		return
	}

	err := h.svc.Enroll(r.Context(), in.CourseID, in.UserID)
	switch {
	case errors.Is(err, course.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, course.ErrNoSeats), errors.Is(err, course.ErrConflict):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case err != nil:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
	default:
		writeJSON(w, http.StatusCreated, map[string]any{
			"courseId": in.CourseID,
			"userId":   in.UserID,
			"status":   "enrolled",
		})
	}
}
