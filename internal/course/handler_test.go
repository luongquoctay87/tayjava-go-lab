package course_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tayjava/go-lab/internal/course"
)

func TestHTTPRoutes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", course.Healthz)
	mux.HandleFunc("GET /api/courses/{id}", course.GetCourse)
	mux.HandleFunc("POST /api/enrollments", course.CreateEnrollment)

	tests := []struct {
		method string
		path   string
		want   int
		body   string
	}{
		{"GET", "/healthz", http.StatusOK, "ok"},
		{"GET", "/api/courses/101", http.StatusOK, "course id=101"},
		{"POST", "/api/enrollments?courseId=101&userId=7", http.StatusCreated, "enrolled user=7 course=101"},
		{"GET", "/api/enrollments", http.StatusMethodNotAllowed, ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != tt.want {
			t.Fatalf("%s %s: status got %d want %d", tt.method, tt.path, rr.Code, tt.want)
		}
		got, _ := io.ReadAll(rr.Body)
		if tt.body != "" && !strings.Contains(string(got), tt.body) {
			t.Fatalf("%s %s: body %q missing %q", tt.method, tt.path, got, tt.body)
		}
	}
}
