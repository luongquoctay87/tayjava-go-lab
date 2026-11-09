package course

import "context"

// Repository is the persistence port used by Service (memory now, Postgres later).
type Repository interface {
	Get(ctx context.Context, id int64) (Course, error)
	Save(ctx context.Context, c Course) (Course, error)
}

// MemoryRepo is an in-process Course store for the API lab.
type MemoryRepo struct {
	data map[int64]Course
	seq  int64
}

// NewMemoryRepo seeds one Course so enroll demos work without the JSON CLI file.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: map[int64]Course{
			1: {ID: 1, Title: "Go Context", Seats: 2, Level: "intermediate"},
		},
		seq: 1,
	}
}

// Get returns a course by id or ErrNotFound.
func (m *MemoryRepo) Get(ctx context.Context, id int64) (Course, error) {
	if err := ctx.Err(); err != nil {
		return Course{}, err
	}
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
