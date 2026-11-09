package enrollment

import "sync/atomic"

// Metrics tracks backlog outcomes. main publishes these through expvar.
type Metrics struct {
	Accepted  atomic.Int64
	Rejected  atomic.Int64
	Processed atomic.Int64
	Failed    atomic.Int64
}
