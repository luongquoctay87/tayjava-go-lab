package course

import (
	"context"
	"sync"
)

// Repository is the persistence port used by Service.
type Repository interface {
	Get(ctx context.Context, id int64) (Course, error)
	Save(ctx context.Context, c Course) (Course, error)
	Enroll(ctx context.Context, courseID, userID int64) error
}

type enrollKey struct {
	courseID int64
	userID   int64
}

// MemoryRepo is an in-process Course store for unit tests and DB-free demos.
type MemoryRepo struct {
	mu       sync.Mutex
	data     map[int64]Course
	enrolled map[enrollKey]struct{}
	seq      int64
}

// NewMemoryRepo seeds one Course so enroll demos work without Postgres.
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

// Enroll takes one seat and records the enrollment under a single lock.
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
