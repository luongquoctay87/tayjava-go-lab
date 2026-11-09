package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tayjava/go-lab/internal/api"
)

func TestHealthzHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	api.Healthz(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body %q", rec.Body.String())
	}
}

func TestProbeHealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(api.Healthz))
	t.Cleanup(srv.Close)
	t.Setenv("HTTP_ADDR", srv.Listener.Addr().String())
	if got := probe(); got != 0 {
		t.Fatalf("probe=%d want 0", got)
	}
}

func TestProbeUnhealthyStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)
	t.Setenv("HTTP_ADDR", srv.Listener.Addr().String())
	if got := probe(); got != 1 {
		t.Fatalf("probe=%d want 1", got)
	}
}

func TestProbeBadAddr(t *testing.T) {
	t.Setenv("HTTP_ADDR", "not-a-hostport")
	if got := probe(); got != 1 {
		t.Fatalf("probe=%d want 1", got)
	}
}
