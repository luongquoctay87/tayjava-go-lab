package enrollment

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/resilience"
)

// Job is an accepted enroll request waiting for a worker.
type Job struct {
	IdempotencyKey string
	UserID         int64
	CourseID       int64
}

var (
	// ErrQueueFull means the bounded channel cannot accept more work.
	ErrQueueFull = errors.New("queue full")
	// ErrShuttingDown means Close has stopped new enqueues.
	ErrShuttingDown = errors.New("shutting down")
	// ErrDuplicate means this user already submitted the same idempotency key.
	ErrDuplicate = errors.New("duplicate idempotency key")
)

// Enroller accepts jobs into a bounded channel and runs workers.
type Enroller struct {
	mu      sync.Mutex
	closed  bool
	seen    map[string]struct{} // accepted userID:key
	jobs    chan Job
	svc     *course.Service
	idem    *resilience.IdempotencyStore
	metrics *Metrics
	wg      sync.WaitGroup
}

// NewEnroller builds a bounded queue (buffer size = queue).
func NewEnroller(svc *course.Service, idem *resilience.IdempotencyStore, m *Metrics, queue int) *Enroller {
	if queue < 1 {
		queue = 1
	}
	if m == nil {
		m = &Metrics{}
	}
	if idem == nil {
		idem = resilience.NewIdempotencyStore()
	}
	return &Enroller{
		seen:    make(map[string]struct{}),
		jobs:    make(chan Job, queue),
		svc:     svc,
		idem:    idem,
		metrics: m,
	}
}

// Enqueue is non-blocking. Duplicate keys are rejected before the send.
// The lock covers the closed flag and the send so Close cannot race a send on a closed channel.
func (e *Enroller) Enqueue(job Job) error {
	k := ""
	if job.IdempotencyKey != "" {
		k = strconv.FormatInt(job.UserID, 10) + ":" + job.IdempotencyKey
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrShuttingDown
	}
	if k != "" {
		if _, dup := e.seen[k]; dup {
			return ErrDuplicate
		}
	}
	select {
	case e.jobs <- job:
		if k != "" {
			e.seen[k] = struct{}{}
		}
		e.metrics.Accepted.Add(1)
		return nil
	default:
		e.metrics.Rejected.Add(1)
		return ErrQueueFull
	}
}

// Depth is the number of jobs still buffered.
func (e *Enroller) Depth() int { return len(e.jobs) }

// Run starts n workers. Each worker uses its own context inside process, never the signal context.
func (e *Enroller) Run(workers int) {
	if workers < 1 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		e.wg.Add(1)
		go e.loop()
	}
}

func (e *Enroller) loop() {
	defer e.wg.Done()
	for job := range e.jobs { // exits after Close, once the buffer is empty
		e.process(job)
	}
}

func (e *Enroller) process(job Job) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if job.IdempotencyKey != "" {
		if _, ok := e.idem.Get(job.IdempotencyKey); ok {
			e.metrics.Processed.Add(1)
			return
		}
	}
	en, err := e.svc.Enroll(ctx, job.CourseID, job.UserID)
	if err != nil {
		e.metrics.Failed.Add(1)
		slog.Error("worker_enroll_failed",
			"course_id", job.CourseID,
			"user_id", job.UserID,
			"idempotency_key", job.IdempotencyKey,
			"err", err.Error(),
		)
		return
	}
	if job.IdempotencyKey != "" {
		e.idem.Put(job.IdempotencyKey, en.ID)
	}
	e.metrics.Processed.Add(1)
	slog.Info("worker_enroll_ok",
		"course_id", job.CourseID,
		"user_id", job.UserID,
		"enrollment_id", en.ID,
		"idempotency_key", job.IdempotencyKey,
	)
}

// Close rejects new enqueues, then closes the job channel. Safe to call more than once.
// It does not wait; the caller bounds Wait with the shutdown deadline.
func (e *Enroller) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.closed {
		e.closed = true
		close(e.jobs)
	}
}

// Wait blocks until workers have drained the closed channel.
func (e *Enroller) Wait() { e.wg.Wait() }
