package enrollment

import (
	"context"
	"errors"
	"log/slog"
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

// ErrQueueFull is returned when the bounded channel cannot accept more work.
var ErrQueueFull = errors.New("enrollment queue full")

// Enroller accepts jobs into a bounded channel and runs workers.
type Enroller struct {
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
		jobs:    make(chan Job, queue),
		svc:     svc,
		idem:    idem,
		metrics: m,
	}
}

// Enqueue non-blocking; returns ErrQueueFull when the buffer is saturated.
func (e *Enroller) Enqueue(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e.jobs <- job:
		e.metrics.Accepted.Add(1)
		e.metrics.setDepth(len(e.jobs))
		return nil
	default:
		e.metrics.Rejected.Add(1)
		return ErrQueueFull
	}
}

// Run starts n workers that drain jobs until Close.
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
	for job := range e.jobs {
		e.metrics.setDepth(len(e.jobs))
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

// Close stops enqueue by closing the jobs channel and waits for workers to drain.
func (e *Enroller) Close() {
	close(e.jobs)
	e.wg.Wait()
}

// Wait blocks until workers exit (after Close).
func (e *Enroller) Wait() { e.wg.Wait() }
