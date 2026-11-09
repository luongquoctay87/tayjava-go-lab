package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Counters is a tiny stdlib-friendly metrics holder for the lab (not Prometheus).
type Counters struct {
	EnrollOK   atomic.Int64
	EnrollFail atomic.Int64
	latencyNS  atomic.Int64
	latencyN   atomic.Int64
}

// ObserveEnroll records an enroll outcome and latency sample.
func (c *Counters) ObserveEnroll(ok bool, d time.Duration) {
	if ok {
		c.EnrollOK.Add(1)
	} else {
		c.EnrollFail.Add(1)
	}
	c.latencyNS.Add(d.Nanoseconds())
	c.latencyN.Add(1)
}

// Handler exposes a temporary text GET /metrics.
func (c *Counters) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := c.latencyN.Load()
		avg := int64(0)
		if n > 0 {
			avg = c.latencyNS.Load() / n
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "enroll_ok_total %d\n", c.EnrollOK.Load())
		_, _ = fmt.Fprintf(w, "enroll_fail_total %d\n", c.EnrollFail.Load())
		_, _ = fmt.Fprintf(w, "enroll_latency_avg_ns %d\n", avg)
	}
}
