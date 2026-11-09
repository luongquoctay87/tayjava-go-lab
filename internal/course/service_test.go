package course_test

import (
	"context"
	"errors"
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
			name:     "two users until full",
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

			if _, err := svc.Enroll(ctx, 10, tt.firstUID); err != nil {
				t.Fatalf("first enroll: %v", err)
			}
			_, err = svc.Enroll(ctx, 10, tt.second.uid)
			if !errors.Is(err, tt.second.wantErr) {
				t.Fatalf("second enroll: got %v want %v", err, tt.second.wantErr)
			}
		})
	}
}

func TestServiceEnrollNotFound(t *testing.T) {
	svc := course.NewService(course.NewMemoryRepo())
	_, err := svc.Enroll(context.Background(), 999, 1)
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

func TestServiceListEnrollments(t *testing.T) {
	ctx := context.Background()
	svc := course.NewService(course.NewMemoryRepo())
	if _, err := svc.Enroll(ctx, 1, 7); err != nil {
		t.Fatal(err)
	}
	list, err := svc.ListEnrollments(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].UserID != 7 {
		t.Fatalf("list=%+v", list)
	}
}
