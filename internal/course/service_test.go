package course_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/tayjava/go-lab/internal/course"
)

func TestServiceEnroll(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		seats    int
		firstUID int64
		second   struct {
			uid     int64
			wantErr error
		}
	}{
		{
			name:     "two users until empty",
			seats:    1,
			firstUID: 7,
			second: struct {
				uid     int64
				wantErr error
			}{uid: 8, wantErr: course.ErrNoSeats},
		},
		{
			name:     "same user conflict",
			seats:    2,
			firstUID: 7,
			second: struct {
				uid     int64
				wantErr error
			}{uid: 7, wantErr: course.ErrConflict},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := course.NewMemoryRepo()
			_, err := repo.Save(ctx, course.Course{ID: 10, Title: "Lab", Seats: tt.seats, Level: "beginner"})
			if err != nil {
				t.Fatal(err)
			}
			svc := course.NewService(repo)

			if err := svc.Enroll(ctx, 10, tt.firstUID); err != nil {
				t.Fatalf("first enroll: %v", err)
			}
			err = svc.Enroll(ctx, 10, tt.second.uid)
			if !errors.Is(err, tt.second.wantErr) {
				t.Fatalf("second enroll: got %v want %v", err, tt.second.wantErr)
			}
		})
	}
}

func TestEnrollConcurrent(t *testing.T) {
	ctx := context.Background()
	repo := course.NewMemoryRepo()
	if _, err := repo.Save(ctx, course.Course{ID: 10, Title: "Lab", Seats: 1, Level: "beginner"}); err != nil {
		t.Fatal(err)
	}
	svc := course.NewService(repo)

	const n = 32
	var wg sync.WaitGroup
	var mu sync.Mutex
	okCount := 0
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(uid int64) {
			defer wg.Done()
			err := svc.Enroll(ctx, 10, uid)
			if err == nil {
				mu.Lock()
				okCount++
				mu.Unlock()
				return
			}
			if !errors.Is(err, course.ErrNoSeats) {
				t.Errorf("uid %d: %v", uid, err)
			}
		}(int64(i + 1))
	}
	wg.Wait()

	if okCount != 1 {
		t.Fatalf("successes=%d want 1", okCount)
	}
	got, err := svc.GetCourse(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if got.Seats != 0 {
		t.Fatalf("seats=%d want 0", got.Seats)
	}
}

func TestServiceEnrollNotFound(t *testing.T) {
	svc := course.NewService(course.NewMemoryRepo())
	err := svc.Enroll(context.Background(), 999, 1)
	if !errors.Is(err, course.ErrNotFound) {
		t.Fatalf("got %v want ErrNotFound", err)
	}
}

func TestServiceCreateAndGet(t *testing.T) {
	ctx := context.Background()
	svc := course.NewService(course.NewMemoryRepo())
	created, err := svc.CreateCourse(ctx, course.Course{Title: "Pointers", Seats: 15, Level: "beginner"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("expected assigned id")
	}
	got, err := svc.GetCourse(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Pointers" || got.Seats != 15 {
		t.Fatalf("got %+v", got)
	}
}
