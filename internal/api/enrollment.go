package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/metrics"
	"github.com/tayjava/go-lab/internal/notify"
	"github.com/tayjava/go-lab/internal/resilience"
)

// EnrollmentHandler exposes enrollment HTTP endpoints.
type EnrollmentHandler struct {
	svc     *course.Service
	notify  *notify.Pool
	metrics *metrics.Counters
	idem    *resilience.IdempotencyStore
}

// NewEnrollmentHandler constructs an EnrollmentHandler.
// notify / metrics / idem may be nil.
func NewEnrollmentHandler(svc *course.Service, pool *notify.Pool, m *metrics.Counters, idem *resilience.IdempotencyStore) EnrollmentHandler {
	return EnrollmentHandler{svc: svc, notify: pool, metrics: m, idem: idem}
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
	h.enroll(w, r, courseID, in.UserID)
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
	slog.Info("list_enrollments",
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
	h.enroll(w, r, in.CourseID, in.UserID)
}

func (h EnrollmentHandler) enroll(w http.ResponseWriter, r *http.Request, courseID, userID int64) {
	key := r.Header.Get("Idempotency-Key")
	slog.Info("enroll_attempt",
		"request_id", RequestIDFrom(r.Context()),
		"course_id", courseID,
		"user_id", userID,
		"idempotency_key", key,
	)
	if key != "" && h.idem != nil {
		if id, ok := h.idem.Get(key); ok {
			WriteJSON(w, http.StatusOK, course.Enrollment{ID: id, CourseID: courseID, UserID: userID})
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	start := time.Now()
	en, err := h.svc.Enroll(ctx, courseID, userID)
	if err != nil {
		h.observe(false, time.Since(start))
		slog.Error("enroll_failed",
			"request_id", RequestIDFrom(r.Context()),
			"course_id", courseID,
			"user_id", userID,
			"err", err.Error(),
		)
		switch {
		case errors.Is(err, course.ErrNoSeats), errors.Is(err, course.ErrConflict), errors.Is(err, course.ErrNotFound):
			MapDomainError(w, err)
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			// A deadline that lands on Commit often surfaces as sql.ErrTxDone.
			// ctx.Err() is the stable signal for a 503 the client may retry.
			w.Header().Set("Retry-After", "1")
			WriteError(w, http.StatusServiceUnavailable, "TIMEOUT", "timeout")
		default:
			MapDomainError(w, err)
		}
		return
	}
	h.observe(true, time.Since(start))
	if key != "" && h.idem != nil {
		h.idem.Put(key, en.ID)
	}
	h.submitNotify(r, en)
	WriteJSON(w, http.StatusCreated, en)
}

func (h EnrollmentHandler) observe(ok bool, d time.Duration) {
	if h.metrics != nil {
		h.metrics.ObserveEnroll(ok, d)
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
