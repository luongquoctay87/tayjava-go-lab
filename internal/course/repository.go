package course

import (
	"context"
	"sync"
)

// Repository is the persistence port used by Service (memory now, Postgres later).
type Repository interface {
	Get(ctx context.Context, id int64) (Course, error)
	Save(ctx context.Context, c Course) (Course, error)
	// Enroll takes one seat and records the enrollment in a single operation.
	Enroll(ctx context.Context, courseID, userID int64) error
}

// enrollKey identifies one user on one course.
type enrollKey struct {
	courseID int64
	userID   int64
}

// MemoryRepo is an in-process Course store for the API lab.
type MemoryRepo struct {
	// mu guards data, enrolled, and seq. Handlers run on many goroutines.
	mu       sync.Mutex
	data     map[int64]Course
	enrolled map[enrollKey]struct{}
	seq      int64
}

// NewMemoryRepo seeds one Course so enroll demos work without the JSON CLI file.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: map[int64]Course{
			1: {ID: 1, Title: "Go Context", Seats: 2, Level: "intermediate"},
		},
		enrolled: make(map[enrollKey]struct{}),
		seq:      1,
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

// Enroll checks seats, decrements one, and records the user under a single lock.
// A failed check writes nothing, so the maps stay consistent.
func (m *MemoryRepo) Enroll(ctx context.Context, courseID, userID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.data[courseID]
	if !ok {
		return ErrNotFound
	}
	if c.Seats <= 0 {
		return ErrNoSeats
	}
	key := enrollKey{courseID: courseID, userID: userID}
	if _, dup := m.enrolled[key]; dup {
		return ErrConflict
	}
	c.Seats--
	m.data[courseID] = c
	m.enrolled[key] = struct{}{}
	return nil
}
