package course

import (
	"context"
	"fmt"
	"sync"
)

// Repository is the persistence port used by Service.
type Repository interface {
	Get(ctx context.Context, id int64) (Course, error)
	Save(ctx context.Context, c Course) (Course, error)
	CountEnrollments(ctx context.Context, courseID int64) (int, error)
	InsertEnrollment(ctx context.Context, courseID, userID int64) (Enrollment, error)
	ListEnrollments(ctx context.Context, courseID int64) ([]Enrollment, error)
}

// MemoryRepo is an in-process Course + Enrollment store for unit tests.
type MemoryRepo struct {
	mu          sync.Mutex
	data        map[int64]Course
	enrollments map[int64]Enrollment
	seq         int64
	enrollSeq   int64
}

// NewMemoryRepo seeds one Course so enroll demos work without Postgres.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: map[int64]Course{
			1: {ID: 1, Title: "Go Context", Seats: 2, Level: "intermediate"},
		},
		enrollments: make(map[int64]Enrollment),
		seq:         1,
	}
}

// Get returns a course by id or ErrNotFound.
func (m *MemoryRepo) Get(ctx context.Context, id int64) (Course, error) {
	if err := ctx.Err(); err != nil {
		return Course{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.data[id]
	if !ok {
		return Course{}, ErrNotFound
	}
	return c, nil
}

// Save inserts or updates a course. ID 0 allocates the next sequence value.
func (m *MemoryRepo) Save(ctx context.Context, c Course) (Course, error) {
	if err := ctx.Err(); err != nil {
		return Course{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID == 0 {
		m.seq++
		c.ID = m.seq
	}
	if c.Level == "" {
		c.Level = "beginner"
	}
	m.data[c.ID] = c
	if c.ID > m.seq {
		m.seq = c.ID
	}
	return c, nil
}

// CountEnrollments returns how many users are enrolled in a course.
func (m *MemoryRepo) CountEnrollments(ctx context.Context, courseID int64) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	for _, e := range m.enrollments {
		if e.CourseID == courseID {
			n++
		}
	}
	return n, nil
}

// InsertEnrollment records a new enrollment or ErrConflict if duplicate.
func (m *MemoryRepo) InsertEnrollment(ctx context.Context, courseID, userID int64) (Enrollment, error) {
	if err := ctx.Err(); err != nil {
		return Enrollment{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%d:%d", courseID, userID)
	for _, e := range m.enrollments {
		if fmt.Sprintf("%d:%d", e.CourseID, e.UserID) == key {
			return Enrollment{}, ErrConflict
		}
	}
	m.enrollSeq++
	e := Enrollment{ID: m.enrollSeq, CourseID: courseID, UserID: userID}
	m.enrollments[e.ID] = e
	return e, nil
}

// ListEnrollments returns all enrollments for a course.
func (m *MemoryRepo) ListEnrollments(ctx context.Context, courseID int64) ([]Enrollment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Enrollment, 0)
	for _, e := range m.enrollments {
		if e.CourseID == courseID {
			out = append(out, e)
		}
	}
	return out, nil
}
