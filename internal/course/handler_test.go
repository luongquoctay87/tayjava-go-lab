package course_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/course"
)

func TestHTTPRoutes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", course.Healthz)
	mux.HandleFunc("GET /api/courses/{id}", course.GetCourse)
	mux.HandleFunc("POST /api/courses", course.CreateCourse)
	mux.HandleFunc("POST /api/enrollments", course.CreateEnrollment)

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)

	t.Run("healthz", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if rr.Code != http.StatusOK || string(mustRead(rr)) != "ok" {
			t.Fatalf("healthz: %d %q", rr.Code, mustRead(rr))
		}
	})

	t.Run("get course json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/courses/101", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d", rr.Code)
		}
		var c course.Course
		if err := json.NewDecoder(rr.Body).Decode(&c); err != nil {
			t.Fatal(err)
		}
		if c.ID != 101 || c.Title == "" {
			t.Fatalf("course %+v", c)
		}
	})

	t.Run("create course json", func(t *testing.T) {
		body := `{"title":"Go Context","seats":30,"level":"intermediate"}`
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/courses", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status %d body %s", rr.Code, mustRead(rr))
		}
		var c course.Course
		if err := json.NewDecoder(rr.Body).Decode(&c); err != nil {
			t.Fatal(err)
		}
		if c.ID != 101 || c.Title != "Go Context" || c.Seats != 30 {
			t.Fatalf("created %+v", c)
		}
	})

	t.Run("enroll", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/enrollments?courseId=101&userId=7", nil))
		if rr.Code != http.StatusCreated {
			t.Fatalf("status %d", rr.Code)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/enrollments", nil))
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status %d want 405", rr.Code)
		}
	})
}

func mustRead(rr *httptest.ResponseRecorder) []byte {
	b, _ := io.ReadAll(rr.Body)
	return b
}
