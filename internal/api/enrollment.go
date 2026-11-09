package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/notify"
)

// EnrollmentHandler exposes enrollment HTTP endpoints.
type EnrollmentHandler struct {
	svc    *course.Service
	notify *notify.Pool
}

// NewEnrollmentHandler constructs an EnrollmentHandler.
// notify may be nil (no post-enroll fan-out).
func NewEnrollmentHandler(svc *course.Service, pool *notify.Pool) EnrollmentHandler {
	return EnrollmentHandler{svc: svc, notify: pool}
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
	h.submitNotify(r, en)
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
	h.submitNotify(r, en)
	WriteJSON(w, http.StatusCreated, en)
}

func (h EnrollmentHandler) submitNotify(r *http.Request, en course.Enrollment) {
	if h.notify == nil {
		return
	}
	// The timeout only bounds the wait for a queue slot. Once queued, workers
	// run the job with the pool context, so a later client disconnect does not
	// drop it.
	ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
	defer cancel()
	if err := h.notify.Submit(ctx, notify.Job{
		CourseID: en.CourseID,
		UserID:   en.UserID,
	}); err != nil {
		log.Printf("notify skipped course=%d user=%d: %v", en.CourseID, en.UserID, err)
	}
}
