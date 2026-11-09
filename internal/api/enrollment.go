package api

import (
	"encoding/json"
	"net/http"

	"github.com/tayjava/go-lab/internal/course"
)

// EnrollmentHandler exposes enrollment HTTP endpoints.
type EnrollmentHandler struct {
	svc *course.Service
}

// NewEnrollmentHandler constructs an EnrollmentHandler.
func NewEnrollmentHandler(svc *course.Service) EnrollmentHandler {
	return EnrollmentHandler{svc: svc}
}

type enrollIn struct {
	CourseID int64 `json:"courseId"`
	UserID   int64 `json:"userId"`
}

func (in enrollIn) validate() string {
	if in.CourseID <= 0 {
		return "courseId must be positive"
	}
	if in.UserID <= 0 {
		return "userId must be positive"
	}
	return ""
}

func decodeEnroll(w http.ResponseWriter, r *http.Request) (enrollIn, string) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in enrollIn
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return in, "invalid json"
	}
	if msg := in.validate(); msg != "" {
		return in, msg
	}
	return in, ""
}

// Create handles POST /api/enrollments with a JSON body.
func (h EnrollmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, msg := decodeEnroll(w, r)
	if msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	if err := h.svc.Enroll(r.Context(), in.CourseID, in.UserID); err != nil {
		MapDomainError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]any{
		"courseId": in.CourseID,
		"userId":   in.UserID,
		"status":   "enrolled",
	})
}
