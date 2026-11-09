package enrollment

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/resilience"
)

func TestEnqueueQueueFull(t *testing.T) {
	svc := course.NewService(course.NewMemoryRepo())
	e := NewEnroller(svc, resilience.NewIdempotencyStore(), &Metrics{}, 1)
	if err := e.Enqueue(context.Background(), Job{CourseID: 1, UserID: 1}); err != nil {
		t.Fatal(err)
	}
	if err := e.Enqueue(context.Background(), Job{CourseID: 1, UserID: 2}); err != ErrQueueFull {
		t.Fatalf("err=%v want ErrQueueFull", err)
	}
}

func TestWorkerProcessesWithIdempotency(t *testing.T) {
	repo := course.NewMemoryRepo()
	_, _ = repo.Save(context.Background(), course.Course{ID: 2, Title: "Adv", Seats: 10})
	svc := course.NewService(repo)
	m := &Metrics{}
	idem := resilience.NewIdempotencyStore()
	e := NewEnroller(svc, idem, m, 8)
	e.Run(2)

	key := "same-key"
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = e.Enqueue(context.Background(), Job{IdempotencyKey: key, CourseID: 2, UserID: 99})
		}()
	}
	wg.Wait()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if m.Processed.Load() >= 1 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	e.Close()

	if m.Processed.Load() < 1 {
		t.Fatalf("processed=%d", m.Processed.Load())
	}
	n, err := repo.CountEnrollments(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("enrollments=%d want 1 (idempotent)", n)
	}
}
