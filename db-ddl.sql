-- TayJava Go lab — Postgres schema for Course + Enrollment (lesson 22+)
-- Column name is `seats` (available capacity), not seats_left.
-- Keep this file in sync with migrations/000001_create_courses and
-- migrations/000002_create_enrollments (golang-migrate up/down).
-- Apply: psql "$DATABASE_URL" -f db-ddl.sql
--    or: migrate -path migrations -database "$DATABASE_URL" up

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
