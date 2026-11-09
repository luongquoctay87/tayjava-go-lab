package enrollment

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tayjava/go-lab/internal/course"
)

func TestCreateAsyncDuplicateAndQueueFull(t *testing.T) {
	e := NewEnroller(course.NewService(course.NewMemoryRepo()), nil, &Metrics{}, 1)
	h := NewHandler(e)

	first := postAsync(h, `{"courseId":1,"userId":7}`, "abc")
	if first.Code != http.StatusAccepted {
		t.Fatalf("first=%d body=%s", first.Code, first.Body.String())
	}
	dup := postAsync(h, `{"courseId":1,"userId":7}`, "abc")
	if dup.Code != http.StatusConflict {
		t.Fatalf("dup=%d body=%s", dup.Code, dup.Body.String())
	}
	full := postAsync(h, `{"courseId":1,"userId":8}`, "other")
	if full.Code != http.StatusServiceUnavailable {
		t.Fatalf("full=%d body=%s", full.Code, full.Body.String())
	}
	if full.Header().Get("Retry-After") != "2" {
		t.Fatalf("Retry-After=%q", full.Header().Get("Retry-After"))
	}
	e.Close()
}

func postAsync(h Handler, body, key string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/enrollments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", key)
	rec := httptest.NewRecorder()
	h.CreateAsync(rec, req)
	return rec
}
