package api_test

import (
	"net/http"
	"net/http/httptest"
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
