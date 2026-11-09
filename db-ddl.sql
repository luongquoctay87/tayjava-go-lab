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

-- Optional (lesson 35+): durable Idempotency-Key store, unique per user.
-- The lab handler uses an in-memory map with the same (user_id, key) scope.
CREATE TABLE IF NOT EXISTS idempotency_keys (
    user_id       BIGINT NOT NULL,
    key           TEXT   NOT NULL,
    enrollment_id BIGINT NOT NULL REFERENCES enrollments (id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, key)
);
