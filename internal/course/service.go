package course

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNotFound means the course id is unknown.
	ErrNotFound = errors.New("course not found")
	// ErrNoSeats means capacity (seats) is exhausted.
	ErrNoSeats = errors.New("no seats left")
	// ErrConflict means the user is already enrolled.
	ErrConflict = errors.New("already enrolled")
	// ErrFull is an alias used in Intermediate lab notes (same as ErrNoSeats).
	ErrFull = ErrNoSeats
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

// Enroll registers userID on courseID when capacity (seats) remains.
// seats is total capacity; enrollments are counted separately (not seats_left).
func (s *Service) Enroll(ctx context.Context, courseID, userID int64) (Enrollment, error) {
	c, err := s.repo.Get(ctx, courseID)
	if err != nil {
		return Enrollment{}, err
	}
	n, err := s.repo.CountEnrollments(ctx, courseID)
	if err != nil {
		return Enrollment{}, err
	}
	if n >= c.Seats {
		return Enrollment{}, ErrNoSeats
	}
	en, err := s.repo.InsertEnrollment(ctx, courseID, userID)
	if err != nil {
		return Enrollment{}, err
	}
	return en, nil
}

// ListEnrollments returns enrollments for a course (404 if course missing).
func (s *Service) ListEnrollments(ctx context.Context, courseID int64) ([]Enrollment, error) {
	if _, err := s.repo.Get(ctx, courseID); err != nil {
		return nil, err
	}
	return s.repo.ListEnrollments(ctx, courseID)
}
