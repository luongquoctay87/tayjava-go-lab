package notify

import (
	"context"
	"errors"
	"sync"
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
		if err := p.Submit(ctx, Job{CourseID: 1, UserID: int64(i)}); err != nil {
			t.Fatalf("submit %d: %v", i, err)
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

	if err := p.Submit(context.Background(), Job{CourseID: 1, UserID: 1}); err != nil {
		t.Fatal(err)
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

func TestSubmitAfterClose(t *testing.T) {
	p := NewPool(context.Background(), 2, 4, func(context.Context, Job) error { return nil })
	p.Close()
	err := p.Submit(context.Background(), Job{CourseID: 1, UserID: 1})
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("got %v want ErrClosed", err)
	}
}

func TestSubmitConcurrentWithClose(t *testing.T) {
	p := NewPool(context.Background(), 2, 1, func(context.Context, Job) error {
		time.Sleep(2 * time.Millisecond)
		return nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = p.Submit(ctx, Job{UserID: int64(id)})
		}(i)
	}
	time.Sleep(5 * time.Millisecond)
	p.Close()
	wg.Wait()
	if err := p.Submit(context.Background(), Job{}); !errors.Is(err, ErrClosed) {
		t.Fatalf("after close: %v", err)
	}
}

func TestCancelDropsQueuedJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	release := make(chan struct{})
	var got atomic.Int64
	p := NewPool(ctx, 1, 8, func(c context.Context, _ Job) error {
		if got.Add(1) == 1 {
			close(started)
			select {
			case <-release:
			case <-c.Done():
			}
		}
		return nil
	})

	for i := 0; i < 5; i++ {
		if err := p.Submit(context.Background(), Job{UserID: int64(i)}); err != nil {
			t.Fatalf("submit: %v", err)
		}
	}
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not start")
	}
	cancel()
	close(release)
	p.Close()
	if got.Load() != 1 {
		t.Fatalf("processed=%d want 1 (queued jobs dropped)", got.Load())
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
