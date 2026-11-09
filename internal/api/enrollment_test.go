package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/course"
)

// stubRepo fakes Repository so handler tests need no Postgres.
type stubRepo struct {
	course      course.Course
	err         error
	enrollCount int
	enrolled    map[string]struct{}
}

func (s stubRepo) Get(ctx context.Context, id int64) (course.Course, error) {
	if err := ctx.Err(); err != nil {
		return course.Course{}, err
	}
	if s.err != nil {
		return course.Course{}, s.err
	}
	if s.course.ID != 0 && id != s.course.ID {
		return course.Course{}, course.ErrNotFound
	}
	return s.course, nil
}

func (s stubRepo) Save(ctx context.Context, c course.Course) (course.Course, error) {
	if err := ctx.Err(); err != nil {
		return course.Course{}, err
	}
	return c, nil
}

func (s stubRepo) CountEnrollments(ctx context.Context, courseID int64) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return s.enrollCount, nil
}

func (s stubRepo) InsertEnrollment(ctx context.Context, courseID, userID int64) (course.Enrollment, error) {
	if err := ctx.Err(); err != nil {
		return course.Enrollment{}, err
	}
	if s.enrolled != nil {
		k := fmt.Sprintf("%d:%d", courseID, userID)
		if _, ok := s.enrolled[k]; ok {
			return course.Enrollment{}, course.ErrConflict
		}
		s.enrolled[k] = struct{}{}
	}
	return course.Enrollment{ID: 1, CourseID: courseID, UserID: userID}, nil
}

func (s stubRepo) ListEnrollments(ctx context.Context, courseID int64) ([]course.Enrollment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func TestEnrollmentCreate(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		want   int
		code   string
		seats0 bool
	}{
		{
			name: "ok",
			body: `{"courseId":1,"userId":7}`,
			want: http.StatusCreated,
		},
		{
			name: "bad json field zero courseId",
			body: `{"courseId":0,"userId":7}`,
			want: http.StatusBadRequest,
			code: "VALIDATION",
		},
		{
			name:   "no seats",
			body:   `{"courseId":1,"userId":8}`,
			want:   http.StatusConflict,
			code:   "NO_SEATS",
			seats0: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := stubRepo{
				course:   course.Course{ID: 1, Title: "Go", Seats: 1, Level: "beginner"},
				enrolled: make(map[string]struct{}),
			}
			if tc.seats0 {
				r.course.Seats = 0
			}
			h := api.NewEnrollmentHandler(course.NewService(r), nil)

			req := httptest.NewRequest(http.MethodPost, "/api/enrollments", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Create(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("status=%d want=%d body=%s", rec.Code, tc.want, rec.Body.String())
			}
			if tc.code == "" {
				return
			}
			var eb api.ErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
				t.Fatal(err)
			}
			if eb.Code != tc.code {
				t.Fatalf("code=%q want=%q", eb.Code, tc.code)
			}
		})
	}
}

func TestEnrollmentCreateNested(t *testing.T) {
	repo := course.NewMemoryRepo()
	svc := course.NewService(repo)
	h := api.NewEnrollmentHandler(svc, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /courses/{id}/enrollments", h.CreateNested)

	req := httptest.NewRequest(http.MethodPost, "/courses/1/enrollments", bytes.NewBufferString(`{"userId":7}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetCourseRoute(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /courses/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") != "42" {
			t.Fatalf("id=%s", r.PathValue("id"))
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/courses/42", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestEnrollmentCreateCanceledContext(t *testing.T) {
	h := api.NewEnrollmentHandler(course.NewService(stubRepo{
		course:   course.Course{ID: 1, Title: "Go", Seats: 1},
		enrolled: make(map[string]struct{}),
	}), nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/enrollments",
		bytes.NewBufferString(`{"courseId":1,"userId":7}`))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want 500 for canceled context", rec.Code)
	}
	var eb api.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Fatal(err)
	}
	if eb.Code != "INTERNAL" {
		t.Fatalf("code=%q want INTERNAL", eb.Code)
	}
}
