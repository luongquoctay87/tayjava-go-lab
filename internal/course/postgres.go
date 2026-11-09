package course

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// PostgresRepo persists Course and Enrollment rows with database/sql.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo wraps an open *sql.DB (already Ping'ed).
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Get loads one course by id.
func (p *PostgresRepo) Get(ctx context.Context, id int64) (Course, error) {
	const q = `SELECT id, title, seats, level FROM courses WHERE id = $1`
	var c Course
	err := p.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Title, &c.Seats, &c.Level)
	if errors.Is(err, sql.ErrNoRows) {
		return Course{}, ErrNotFound
	}
	return c, err
}

// Save inserts (id == 0) or updates a course. Uses $n placeholders — never string-concat SQL.
func (p *PostgresRepo) Save(ctx context.Context, c Course) (Course, error) {
	if c.Level == "" {
		c.Level = "beginner"
	}
	if c.ID == 0 {
		const q = `INSERT INTO courses (title, seats, level) VALUES ($1, $2, $3) RETURNING id`
		err := p.db.QueryRowContext(ctx, q, c.Title, c.Seats, c.Level).Scan(&c.ID)
		return c, err
	}
	const q = `UPDATE courses SET title = $1, seats = $2, level = $3 WHERE id = $4`
	res, err := p.db.ExecContext(ctx, q, c.Title, c.Seats, c.Level, c.ID)
	if err != nil {
		return Course{}, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return Course{}, err
	}
	if n == 0 {
		return Course{}, ErrNotFound
	}
	return c, nil
}

// CountEnrollments returns enrollment rows for a course.
func (p *PostgresRepo) CountEnrollments(ctx context.Context, courseID int64) (int, error) {
	const q = `SELECT COUNT(*) FROM enrollments WHERE course_id = $1`
	var n int
	err := p.db.QueryRowContext(ctx, q, courseID).Scan(&n)
	return n, err
}

// InsertEnrollment inserts a row; unique (course_id, user_id) maps to ErrConflict.
func (p *PostgresRepo) InsertEnrollment(ctx context.Context, courseID, userID int64) (Enrollment, error) {
	const q = `
INSERT INTO enrollments (course_id, user_id) VALUES ($1, $2)
RETURNING id, course_id, user_id`
	var e Enrollment
	err := p.db.QueryRowContext(ctx, q, courseID, userID).Scan(&e.ID, &e.CourseID, &e.UserID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return Enrollment{}, ErrConflict
		}
		return Enrollment{}, err
	}
	return e, nil
}

// ListEnrollments returns enrollments for a course ordered by id.
func (p *PostgresRepo) ListEnrollments(ctx context.Context, courseID int64) ([]Enrollment, error) {
	const q = `SELECT id, course_id, user_id FROM enrollments WHERE course_id = $1 ORDER BY id`
	rows, err := p.db.QueryContext(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Enrollment, 0)
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.ID, &e.CourseID, &e.UserID); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// OpenDB opens a Postgres pool and verifies connectivity with PingContext.
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// OpenDBWithRetry retries OpenDB so Compose can wait for Postgres readiness.
func OpenDBWithRetry(dsn string, attempts int, delay time.Duration) (*sql.DB, error) {
	var last error
	for i := 0; i < attempts; i++ {
		db, err := OpenDB(dsn)
		if err == nil {
			return db, nil
		}
		last = err
		time.Sleep(delay)
	}
	return nil, last
}
