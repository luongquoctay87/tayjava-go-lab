package metrics

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCountersHandler(t *testing.T) {
	c := &Counters{}
	c.ObserveEnroll(true, 10*time.Millisecond)
	c.ObserveEnroll(false, 5*time.Millisecond)
	rr := httptest.NewRecorder()
	c.Handler()(rr, httptest.NewRequest("GET", "/metrics", nil))
	body := rr.Body.String()
	if !strings.Contains(body, "enroll_ok_total 1") || !strings.Contains(body, "enroll_fail_total 1") {
		t.Fatalf("body=%q", body)
	}
}
