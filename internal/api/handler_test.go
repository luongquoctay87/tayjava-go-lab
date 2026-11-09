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
	courses := api.NewCourseHandler(svc)
	enrollments := api.NewEnrollmentHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Healthz)
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
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/courses/1", nil))
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
		body1 := `{"courseId":1,"userId":7}`
		rr1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodPost, "/api/enrollments", strings.NewReader(body1))
		h.ServeHTTP(rr1, req1)
		if rr1.Code != http.StatusCreated {
			t.Fatalf("first enroll status %d", rr1.Code)
		}

		body2 := `{"courseId":1,"userId":8}`
		rr2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodPost, "/api/enrollments", strings.NewReader(body2))
		h.ServeHTTP(rr2, req2)
		if rr2.Code != http.StatusCreated {
			t.Fatalf("second enroll status %d", rr2.Code)
		}

		body3 := `{"courseId":1,"userId":9}`
		rr3 := httptest.NewRecorder()
		req3 := httptest.NewRequest(http.MethodPost, "/api/enrollments", strings.NewReader(body3))
		h.ServeHTTP(rr3, req3)
		if rr3.Code != http.StatusConflict {
			t.Fatalf("third enroll status %d want 409", rr3.Code)
		}
	})

	t.Run("create course", func(t *testing.T) {
		body := `{"title":"Pointers","seats":10,"level":"beginner"}`
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/courses", strings.NewReader(body))
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status %d", rr.Code)
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
