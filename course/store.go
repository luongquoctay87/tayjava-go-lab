package course

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned when a course ID is missing from a store.
var ErrNotFound = errors.New("course: not found")

// CourseStore is the small persistence surface used by Enroll.
type CourseStore interface {
	FindByID(id int64) (Course, error)
	Save(c Course) error
}

// MemoryStore keeps courses in a map for demos and tests.
type MemoryStore struct {
	items map[int64]Course
}

// NewMemoryStore creates an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[int64]Course)}
}

// FindByID returns ErrNotFound when the id is absent.
func (s *MemoryStore) FindByID(id int64) (Course, error) {
	c, ok := s.items[id]
	if !ok {
		return Course{}, ErrNotFound
	}
	return c, nil
}

// Save upserts a course by ID.
func (s *MemoryStore) Save(c Course) error {
	if s.items == nil {
		s.items = make(map[int64]Course)
	}
	s.items[c.ID] = c
	return nil
}

// Enroll looks up a course and wraps lookup failures with %w.
func Enroll(store CourseStore, courseID, userID int64) error {
	_, err := store.FindByID(courseID)
	if err != nil {
		return fmt.Errorf("enroll user %d into course %d: %w", userID, courseID, err)
	}
	return nil
}
