package course

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNotFound means the course id is unknown.
	ErrNotFound = errors.New("course not found")
	// ErrNoSeats means capacity is exhausted.
	ErrNoSeats = errors.New("no seats left")
	// ErrConflict means the user is already enrolled.
	ErrConflict = errors.New("already enrolled")
)

// Service holds Course use-cases. It does not import net/http.
type Service struct {
	repo Repository
}

// NewService wires a repository implementation.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetCourse loads a course by id.
func (s *Service) GetCourse(ctx context.Context, id int64) (Course, error) {
	return s.repo.Get(ctx, id)
}

// CreateCourse persists a new course (id assigned by repository when 0).
func (s *Service) CreateCourse(ctx context.Context, c Course) (Course, error) {
	if c.Title == "" {
		return Course{}, fmt.Errorf("title required")
	}
	if c.Seats < 0 {
		return Course{}, fmt.Errorf("seats must not be negative")
	}
	c.ID = 0
	return s.repo.Save(ctx, c)
}

// Enroll reserves one seat for userID on courseID.
// Only ErrNotFound from Get stays a 404; other Get errors are wrapped so they become 500.
func (s *Service) Enroll(ctx context.Context, courseID, userID int64) error {
	if _, err := s.repo.Get(ctx, courseID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return err
		}
		return fmt.Errorf("get course %d: %w", courseID, err)
	}
	return s.repo.Enroll(ctx, courseID, userID)
}
