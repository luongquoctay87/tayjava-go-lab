package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tayjava/go-lab/internal/course"
)

// CourseHandler exposes Course HTTP endpoints.
type CourseHandler struct {
	svc      *course.Service
	seatCap  int
}

// NewCourseHandler constructs a CourseHandler.
// seatCapDefault is applied when the create body omits seats / sends 0.
func NewCourseHandler(svc *course.Service, seatCapDefault int) CourseHandler {
	if seatCapDefault < 1 {
		seatCapDefault = 30
	}
	return CourseHandler{svc: svc, seatCap: seatCapDefault}
}

type createCourseIn struct {
	Title string `json:"title"`
	Seats int    `json:"seats"`
	Level string `json:"level"`
}

func (in createCourseIn) validate() string {
	if strings.TrimSpace(in.Title) == "" {
		return "title is required"
	}
	if in.Seats < 0 {
		return "seats must not be negative"
	}
	return ""
}

func decodeCreateCourse(w http.ResponseWriter, r *http.Request) (createCourseIn, string) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in createCourseIn
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

// Get handles GET /api/courses/{id}.
func (h CourseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "VALIDATION", "id must be a positive integer")
		return
	}
	c, err := h.svc.GetCourse(r.Context(), id)
	if err != nil {
		MapDomainError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, c)
}

// Create handles POST /api/courses.
func (h CourseHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, msg := decodeCreateCourse(w, r)
	if msg != "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION", msg)
		return
	}
	seats := in.Seats
	if seats == 0 {
		seats = h.seatCap
	}
	created, err := h.svc.CreateCourse(r.Context(), course.Course{
		Title: in.Title,
		Seats: seats,
		Level: in.Level,
	})
	if err != nil {
		MapDomainError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

// Healthz handles GET /healthz.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
