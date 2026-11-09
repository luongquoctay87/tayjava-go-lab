package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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
	idem    *resilience.Idempotency
}

// NewEnrollmentHandler constructs an EnrollmentHandler.
// notify / metrics / idem may be nil. A non-nil idem store requires Idempotency-Key.
func NewEnrollmentHandler(svc *course.Service, pool *notify.Pool, m *metrics.Counters, idem *resilience.Idempotency) EnrollmentHandler {
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

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
}

func decodeEnrollBody(raw []byte) (enrollBody, string) {
	var in enrollBody
	dec := json.NewDecoder(bytes.NewReader(raw))
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
	raw, err := readBody(w, r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION", "invalid body")
		return
	}
	in, msg := decodeEnrollBody(raw)
	if msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	if msg := in.validateNested(); msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	h.enroll(w, r, courseID, in.UserID, raw)
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
	raw, err := readBody(w, r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION", "invalid body")
		return
	}
	in, msg := decodeEnrollBody(raw)
	if msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	if msg := in.validateFlat(); msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	h.enroll(w, r, in.CourseID, in.UserID, raw)
}

func (h EnrollmentHandler) enroll(w http.ResponseWriter, r *http.Request, courseID, userID int64, rawBody []byte) {
	key := r.Header.Get("Idempotency-Key")
	slog.Info("enroll_attempt",
		"request_id", RequestIDFrom(r.Context()),
		"course_id", courseID,
		"user_id", userID,
		"idempotency_key", key,
	)
	if h.idem == nil {
		status, body := h.performEnroll(r, courseID, userID)
		writeStored(w, status, body)
		return
	}
	if key == "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", "missing Idempotency-Key")
		return
	}

	status, resp, err := h.idem.Do(userID, key, rawBody, func() (int, []byte) {
		return h.performEnroll(r, courseID, userID)
	})
	switch {
	case errors.Is(err, resilience.ErrKeyReused):
		WriteError(w, http.StatusUnprocessableEntity, "KEY_REUSED", err.Error())
	case errors.Is(err, resilience.ErrInProgress):
		w.Header().Set("Retry-After", "1")
		WriteError(w, http.StatusConflict, "IN_PROGRESS", err.Error())
	default:
		writeStored(w, status, resp)
	}
}

// performEnroll reserves a seat. It runs outside the idempotency mutex.
func (h EnrollmentHandler) performEnroll(r *http.Request, courseID, userID int64) (int, []byte) {
	start := time.Now()
	en, err := h.svc.Enroll(r.Context(), courseID, userID)
	if err != nil {
		h.observe(false, time.Since(start))
		slog.Error("enroll_failed",
			"request_id", RequestIDFrom(r.Context()),
			"course_id", courseID,
			"user_id", userID,
			"err", err.Error(),
		)
		status, body := domainErrorResponse(err)
		if status >= 500 {
			slog.Error("internal_error", "err", err.Error())
		}
		return status, jsonBody(body)
	}
	h.observe(true, time.Since(start))
	h.submitNotify(r, en)
	return http.StatusCreated, jsonBody(en)
}

func writeStored(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
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
