CREATE TABLE IF NOT EXISTS courses (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    seats      INT  NOT NULL CHECK (seats >= 0),
    level      TEXT NOT NULL DEFAULT 'beginner',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
