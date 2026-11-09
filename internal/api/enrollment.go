package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

type enrollBody struct {
	UserID   int64 `json:"userId"`
	CourseID int64 `json:"courseId"` // used by flat /api/enrollments
}

func (in enrollBody) validateNested() string {
	if in.UserID <= 0 {
		return "userId must be positive"
	}
	return ""
}

func (in enrollBody) validateFlat() string {
	if in.CourseID <= 0 {
		return "courseId must be positive"
	}
	if in.UserID <= 0 {
		return "userId must be positive"
	}
	return ""
}

func decodeEnrollBody(w http.ResponseWriter, r *http.Request) (enrollBody, string) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in enrollBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return in, "invalid json"
	}
	return in, ""
}

// CreateNested handles POST /courses/{id}/enrollments.
func (h EnrollmentHandler) CreateNested(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || courseID <= 0 {
		WriteError(w, http.StatusBadRequest, "VALIDATION", "id must be a positive integer")
		return
	}
	in, msg := decodeEnrollBody(w, r)
	if msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	if msg := in.validateNested(); msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	en, err := h.svc.Enroll(r.Context(), courseID, in.UserID)
	if err != nil {
		MapDomainError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, en)
}

// List handles GET /courses/{id}/enrollments.
func (h EnrollmentHandler) List(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || courseID <= 0 {
		WriteError(w, http.StatusBadRequest, "VALIDATION", "id must be a positive integer")
		return
	}
	list, err := h.svc.ListEnrollments(r.Context(), courseID)
	if err != nil {
		MapDomainError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

// Create handles POST /api/enrollments with courseId + userId (kept for earlier lab curls).
func (h EnrollmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, msg := decodeEnrollBody(w, r)
	if msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	if msg := in.validateFlat(); msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	en, err := h.svc.Enroll(r.Context(), in.CourseID, in.UserID)
	if err != nil {
		MapDomainError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, en)
}
