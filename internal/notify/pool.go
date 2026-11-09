package notify

import (
	"context"
	"log"
	"sync"
)

// Job is a lightweight post-enrollment notify unit of work (fake SMTP in the lab).
type Job struct {
	CourseID int64
	UserID   int64
}

// Pool is a fixed-size worker pool with a bounded job queue (fan-out from one channel).
type Pool struct {
	jobs chan Job
	wg   sync.WaitGroup
}

// NewPool starts workers goroutines. Pass context.Background() when you still need
// to drain via Close after HTTP Shutdown (do not share the signal ctx — see lesson 29).
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
	p := &Pool{jobs: make(chan Job, queue)}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
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

// Submit enqueues a job. Returns false if ctx is canceled before the job is accepted.
// A full buffer blocks until space is available or ctx cancels.
func (p *Pool) Submit(ctx context.Context, job Job) bool {
	select {
	case <-ctx.Done():
		return false
	case p.jobs <- job:
		return true
	}
}

// Close stops accepting jobs and waits for workers to finish draining the queue.
func (p *Pool) Close() {
	close(p.jobs)
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
