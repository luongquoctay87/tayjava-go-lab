package notify

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolProcessesAllJobs(t *testing.T) {
	var got atomic.Int64
	ctx := context.Background()
	p := NewPool(ctx, 4, 64, func(_ context.Context, job Job) error {
		got.Add(1)
		_ = job
		return nil
	})

	const n = 100
	for i := 0; i < n; i++ {
		if !p.Submit(ctx, Job{CourseID: 1, UserID: int64(i)}) {
			t.Fatalf("submit %d rejected", i)
		}
	}
	p.Close()

	if got.Load() != n {
		t.Fatalf("processed=%d want=%d", got.Load(), n)
	}
}

func TestPoolStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	p := NewPool(ctx, 2, 8, func(c context.Context, _ Job) error {
		select {
		case <-started:
		default:
			close(started)
		}
		<-c.Done()
		return c.Err()
	})

	if !p.Submit(context.Background(), Job{CourseID: 1, UserID: 1}) {
		t.Fatal("submit failed")
	}
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not start")
	}
	cancel()

	done := make(chan struct{})
	go func() {
		p.Close()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Close hung after cancel")
	}
}

func TestFanInAggregates(t *testing.T) {
	ctx := context.Background()
	a := make(chan int64, 2)
	b := make(chan int64, 2)
	a <- 1
	a <- 2
	b <- 3
	close(a)
	close(b)

	out := FanIn(ctx, a, b)
	sum := int64(0)
	for v := range out {
		sum += v
	}
	if sum != 6 {
		t.Fatalf("sum=%d want=6", sum)
	}
}
