package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tayjava/go-lab/internal/api"
)

func TestRecoveryCatchesPanic(t *testing.T) {
	h := api.Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status got %d want 500", rr.Code)
	}
}

func TestLoggingRecordsStatus(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	h := api.Logging(inner)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/courses", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status got %d want 201", rr.Code)
	}
}

func TestRequestIDValidatedOrGenerated(t *testing.T) {
	h := api.Logging(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "req-abc-123")
	h.ServeHTTP(rr, req)
	if got := rr.Header().Get("X-Request-ID"); got != "req-abc-123" {
		t.Fatalf("valid id got %q", got)
	}

	for _, bad := range []string{"", "has space", "bad\nid", strings.Repeat("a", 500)} {
		rr = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", bad)
		h.ServeHTTP(rr, req)
		got := rr.Header().Get("X-Request-ID")
		if got == "" || got == bad || len(got) > 64 || !requestIDCharset(got) {
			t.Fatalf("bad %q replaced with %q", bad, got)
		}
	}
}

func requestIDCharset(s string) bool {
	for _, c := range s {
		ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_'
		if !ok {
			return false
		}
	}
	return s != ""
}
