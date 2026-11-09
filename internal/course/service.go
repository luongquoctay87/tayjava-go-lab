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
	repo     Repository
	enrolled map[string]struct{}
}

// NewService wires a repository implementation.
func NewService(repo Repository) *Service {
	return &Service{repo: repo, enrolled: make(map[string]struct{})}
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
func (s *Service) Enroll(ctx context.Context, courseID, userID int64) error {
	c, err := s.repo.Get(ctx, courseID)
	if err != nil {
		return err
	}
	if c.Seats <= 0 {
		return ErrNoSeats
	}
	key := fmt.Sprintf("%d:%d", courseID, userID)
	if _, ok := s.enrolled[key]; ok {
		return ErrConflict
	}
	s.enrolled[key] = struct{}{}
	c.Seats--
	_, err = s.repo.Save(ctx, c)
	return err
}
