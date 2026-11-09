-- TayJava Go lab — Postgres schema for Course + Enrollment (lesson 22+)
-- Column name is `seats` (total capacity; not seats_left).
-- Apply: psql "$DATABASE_URL" -f db-ddl.sql

CREATE TABLE IF NOT EXISTS courses (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    seats      INT  NOT NULL CHECK (seats >= 0),
    level      TEXT NOT NULL DEFAULT 'beginner',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS enrollments (
    id         BIGSERIAL PRIMARY KEY,
    course_id  BIGINT NOT NULL REFERENCES courses (id),
    user_id    BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, user_id)
);

-- Optional (lesson 35+): durable Idempotency-Key store for enroll retries.
-- The in-process lab also keeps an in-memory map; prefer this table in Postgres.
CREATE TABLE IF NOT EXISTS idempotency_keys (
    key            TEXT PRIMARY KEY,
    enrollment_id  BIGINT NOT NULL REFERENCES enrollments (id),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
