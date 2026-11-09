package enrollment

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Metrics tracks backlog and processing outcomes for the advanced lab.
type Metrics struct {
	Accepted  atomic.Int64
	Rejected  atomic.Int64
	Processed atomic.Int64
	Failed    atomic.Int64
	depth     atomic.Int64
}

func (m *Metrics) setDepth(n int) { m.depth.Store(int64(n)) }

// Handler exposes text metrics.
func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "enroll_accepted_total %d\n", m.Accepted.Load())
		_, _ = fmt.Fprintf(w, "enroll_rejected_total %d\n", m.Rejected.Load())
		_, _ = fmt.Fprintf(w, "enroll_processed_total %d\n", m.Processed.Load())
		_, _ = fmt.Fprintf(w, "enroll_failed_total %d\n", m.Failed.Load())
		_, _ = fmt.Fprintf(w, "enroll_queue_depth %d\n", m.depth.Load())
	}
}
