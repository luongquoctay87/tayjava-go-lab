package course_test

import (
	"errors"
	"testing"

	"github.com/tayjava/go-lab/course"
)

func TestRemainingSeatsOK(t *testing.T) {
	got, err := course.RemainingSeats(30, 10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != 20 {
		t.Fatalf("got %d, want 20", got)
	}
}

func TestRemainingSeats(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		enrolled int
		want     int
		wantErr  error
	}{
		{name: "normal", capacity: 30, enrolled: 10, want: 20},
		{name: "full", capacity: 10, enrolled: 10, want: 0},
		{name: "over", capacity: 5, enrolled: 6, wantErr: course.ErrInvalidEnrollment},
		{name: "negative", capacity: -1, enrolled: 0, wantErr: errors.New("negative seats")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := course.RemainingSeats(tt.capacity, tt.enrolled)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("want error %v, got nil", tt.wantErr)
				}
				if tt.wantErr == course.ErrInvalidEnrollment {
					if !errors.Is(err, course.ErrInvalidEnrollment) {
						t.Fatalf("want ErrInvalidEnrollment, got %v", err)
					}
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("got err %q, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFindByIDNotFound(t *testing.T) {
	store := course.NewMemoryStore()
	_, err := store.FindByID(42)
	if !errors.Is(err, course.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
