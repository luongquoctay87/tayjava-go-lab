package course

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// pgUniqueViolation is the Postgres SQLSTATE for unique_violation.
const pgUniqueViolation = "23505"

// PostgresRepo persists Course rows with database/sql.
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

// Enroll decrements one seat and inserts an enrollment row.
// The two statements are not wrapped in a transaction; lesson 38 does that.
// If the insert fails after the update, the seat is already gone.
func (p *PostgresRepo) Enroll(ctx context.Context, courseID, userID int64) error {
	// seats > 0 is in the same UPDATE, so two concurrent requests cannot both take the last seat.
	res, err := p.db.ExecContext(ctx,
		`UPDATE courses SET seats = seats - 1 WHERE id = $1 AND seats > 0`, courseID)
	if err != nil {
		return fmt.Errorf("reserve seat: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("reserve seat: %w", err)
	}
	if n == 0 {
		return ErrNoSeats
	}

	_, err = p.db.ExecContext(ctx,
		`INSERT INTO enrollments (course_id, user_id) VALUES ($1, $2)`, courseID, userID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return ErrConflict
	}
	if err != nil {
		return fmt.Errorf("insert enrollment: %w", err)
	}
	return nil
}

// OpenDB opens a Postgres pool and verifies connectivity with PingContext.
// Callers pass DATABASE_URL. Idle is set equal to open so the pool
// does not close connections down to the database/sql default of 2.
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
