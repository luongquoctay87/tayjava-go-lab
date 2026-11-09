package notify

import (
	"context"
	"errors"
	"log"
	"sync"
)

// ErrClosed is returned by Submit after Close has started.
var ErrClosed = errors.New("notify: pool closed")

// Job is a lightweight post-enrollment notify unit of work (fake SMTP in the lab).
type Job struct {
	CourseID int64
	UserID   int64
}

// Pool is a fixed-size worker pool with a bounded job queue (fan-out from one channel).
type Pool struct {
	jobs      chan Job
	done      chan struct{} // closed when Close starts: stop accepting jobs
	mu        sync.RWMutex  // Submit holds RLock while sending; Close holds Lock for close(jobs)
	closeOnce sync.Once
	wg        sync.WaitGroup
}

// NewPool starts workers goroutines.
// ctx is the process context: canceling it stops workers after the job in
// progress and drops jobs still queued. Close is different: it stops accepting
// new jobs and workers drain the queue before exiting.
func NewPool(ctx context.Context, workers, queue int, handle func(context.Context, Job) error) *Pool {
	if workers < 1 {
		workers = 1
	}
	if queue < 1 {
		queue = 1
	}
	if handle == nil {
		handle = FakeSMTP
	}
	p := &Pool{
		jobs: make(chan Job, queue),
		done: make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				// Prefer cancel over a non-empty queue so leftover jobs are dropped.
				if ctx.Err() != nil {
					return
				}
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return // Close closed jobs and the queue is empty
					}
					if err := handle(ctx, job); err != nil {
						log.Printf("notify job failed: %v", err)
					}
				}
			}
		}()
	}
	return p
}

// Submit waits at most until ctx is done. It returns ErrClosed after Close
// and never panics on a closed jobs channel.
func (p *Pool) Submit(ctx context.Context, job Job) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// A send on a closed channel panics even inside select, so check first.
	select {
	case <-p.done:
		return ErrClosed
	default:
	}

	select {
	case p.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.done:
		return ErrClosed
	}
}

// Close stops accepting jobs, drains the queue, and waits for workers.
// It is safe to call more than once and concurrently with Submit.
func (p *Pool) Close() {
	p.closeOnce.Do(func() {
		close(p.done) // wake any Submit blocked on a full queue so it releases RLock
		p.mu.Lock()
		close(p.jobs)
		p.mu.Unlock()
	})
	p.wg.Wait()
}

// FakeSMTP pretends to send an enrollment email.
func FakeSMTP(_ context.Context, job Job) error {
	log.Printf("enrolled user=%d course=%d", job.UserID, job.CourseID)
	return nil
}

// FanIn merges multiple result channels into one (lab fan-in demo).
// The returned channel is closed when every input is drained.
func FanIn(ctx context.Context, channels ...<-chan int64) <-chan int64 {
	out := make(chan int64)
	var wg sync.WaitGroup
	forward := func(c <-chan int64) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go forward(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
