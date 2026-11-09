package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/course"
)

func newTestHandler() http.Handler {
	repo := course.NewMemoryRepo()
	svc := course.NewService(repo)
	courses := api.NewCourseHandler(svc, 30)
	enrollments := api.NewEnrollmentHandler(svc, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Healthz)
	mux.HandleFunc("GET /courses/{id}", courses.Get)
	mux.HandleFunc("POST /courses", courses.Create)
	mux.HandleFunc("POST /courses/{id}/enrollments", enrollments.CreateNested)
	mux.HandleFunc("GET /courses/{id}/enrollments", enrollments.List)
	mux.HandleFunc("GET /api/courses/{id}", courses.Get)
	mux.HandleFunc("POST /api/courses", courses.Create)
	mux.HandleFunc("POST /api/enrollments", enrollments.Create)

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)
	return h
}

func TestAPILayers(t *testing.T) {
	h := newTestHandler()

	t.Run("get seeded course", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/courses/1", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d", rr.Code)
		}
		var c course.Course
		if err := json.NewDecoder(rr.Body).Decode(&c); err != nil {
			t.Fatal(err)
		}
		if c.ID != 1 || c.Seats != 2 {
			t.Fatalf("course %+v", c)
		}
	})

	t.Run("enroll until no seats", func(t *testing.T) {
		rr1 := httptest.NewRecorder()
		h.ServeHTTP(rr1, httptest.NewRequest(http.MethodPost, "/courses/1/enrollments", strings.NewReader(`{"userId":7}`)))
		if rr1.Code != http.StatusCreated {
			t.Fatalf("first enroll status %d body=%s", rr1.Code, rr1.Body.String())
		}

		rr2 := httptest.NewRecorder()
		h.ServeHTTP(rr2, httptest.NewRequest(http.MethodPost, "/courses/1/enrollments", strings.NewReader(`{"userId":8}`)))
		if rr2.Code != http.StatusCreated {
			t.Fatalf("second enroll status %d", rr2.Code)
		}

		rr3 := httptest.NewRecorder()
		h.ServeHTTP(rr3, httptest.NewRequest(http.MethodPost, "/courses/1/enrollments", strings.NewReader(`{"userId":9}`)))
		if rr3.Code != http.StatusConflict {
			t.Fatalf("third enroll status %d want 409", rr3.Code)
		}
		var eb api.ErrorBody
		if err := json.Unmarshal(rr3.Body.Bytes(), &eb); err != nil {
			t.Fatal(err)
		}
		if eb.Code != "NO_SEATS" {
			t.Fatalf("code=%q want NO_SEATS", eb.Code)
		}
	})

	t.Run("validation userId zero", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/courses/1/enrollments", strings.NewReader(`{"userId":0}`)))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status %d want 400", rr.Code)
		}
		var eb api.ErrorBody
		if err := json.Unmarshal(rr.Body.Bytes(), &eb); err != nil {
			t.Fatal(err)
		}
		if eb.Code != "VALIDATION" {
			t.Fatalf("code=%q want VALIDATION", eb.Code)
		}
	})

	t.Run("create course", func(t *testing.T) {
		body := `{"title":"Pointers","seats":10,"level":"beginner"}`
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/courses", strings.NewReader(body))
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status %d body=%s", rr.Code, rr.Body.String())
		}
		var c course.Course
		if err := json.NewDecoder(rr.Body).Decode(&c); err != nil {
			t.Fatal(err)
		}
		if c.ID == 0 || c.Title != "Pointers" {
			t.Fatalf("created %+v", c)
		}
	})
}
