package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/notify"
)

// EnrollmentHandler exposes enrollment HTTP endpoints.
type EnrollmentHandler struct {
	svc    *course.Service
	notify *notify.Pool
	logger *slog.Logger
}

// NewEnrollmentHandler constructs an EnrollmentHandler.
// notify may be nil (no post-enroll fan-out).
// logger is the process logger created after slog.SetDefault; the child logger
// with component=enrollment is built here, not in a package-level variable.
func NewEnrollmentHandler(svc *course.Service, pool *notify.Pool, logger *slog.Logger) EnrollmentHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return EnrollmentHandler{
		svc:    svc,
		notify: pool,
		logger: logger.With("component", "enrollment"),
	}
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
	h.logAttempt(r.Context(), courseID, in.UserID)
	en, err := h.svc.Enroll(r.Context(), courseID, in.UserID)
	if err != nil {
		h.logEnrollErr(r.Context(), courseID, in.UserID, err)
		MapDomainError(w, err)
		return
	}
	h.logEnrollOK(r.Context(), courseID, in.UserID)
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
	h.logger.Info("list_enrollments",
		"request_id", RequestIDFrom(r.Context()),
		"course_id", courseID,
		"count", len(list),
	)
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
	h.logAttempt(r.Context(), in.CourseID, in.UserID)
	en, err := h.svc.Enroll(r.Context(), in.CourseID, in.UserID)
	if err != nil {
		h.logEnrollErr(r.Context(), in.CourseID, in.UserID, err)
		MapDomainError(w, err)
		return
	}
	h.logEnrollOK(r.Context(), in.CourseID, in.UserID)
	h.submitNotify(r, en)
	WriteJSON(w, http.StatusCreated, en)
}

func (h EnrollmentHandler) logAttempt(ctx context.Context, courseID, userID int64) {
	h.logger.Info("enrollment attempt",
		"request_id", RequestIDFrom(ctx),
		"course_id", courseID,
		"user_id", userID,
	)
}

func (h EnrollmentHandler) logEnrollOK(ctx context.Context, courseID, userID int64) {
	h.logger.Info("enrollment ok",
		"request_id", RequestIDFrom(ctx),
		"course_id", courseID,
		"user_id", userID,
	)
}

// logEnrollErr records business rejections at Info. Unexpected failures are Error.
func (h EnrollmentHandler) logEnrollErr(ctx context.Context, courseID, userID int64, err error) {
	rid := RequestIDFrom(ctx)
	switch {
	case errors.Is(err, course.ErrNoSeats):
		h.logger.Info("enrollment rejected",
			"request_id", rid,
			"course_id", courseID,
			"user_id", userID,
			"reason", "no_seats",
		)
	case errors.Is(err, course.ErrConflict):
		h.logger.Info("enrollment rejected",
			"request_id", rid,
			"course_id", courseID,
			"user_id", userID,
			"reason", "conflict",
		)
	default:
		h.logger.Error("enrollment failed",
			"request_id", rid,
			"course_id", courseID,
			"user_id", userID,
			"err", err,
		)
	}
}

func (h EnrollmentHandler) submitNotify(r *http.Request, en course.Enrollment) {
	if h.notify == nil {
		return
	}
	// Use Background so client cancel does not drop accepted notify work.
	_ = h.notify.Submit(context.Background(), notify.Job{
		CourseID: en.CourseID,
		UserID:   en.UserID,
	})
	_ = r
}
