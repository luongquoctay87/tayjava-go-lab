# Lesson 25 — Env config / twelve-factor light

Branch: `lesson-25`  
Article: [Config bằng env và twelve-factor nhẹ trong Golang](https://tayjava.com/vi/technologies/config-bang-env-va-twelve-factor-nhe-trong-golang)

## Goal of this branch

Load a `Config` struct from the environment (`HTTP_ADDR`, `DATABASE_URL`, `ENROLLMENT_CAP_DEFAULT`, `SHUTDOWN_TIMEOUT`), fail fast when required vars are missing, and keep secrets out of Git via `.env.example` + gitignored `.env`.

## Changes vs `lesson-24`

| Add / change | Why |
|--------------|-----|
| `internal/config/` | `Load()` + unit tests |
| `cmd/api/main.go` | Wire config; Postgres required |
| `internal/api/course.go` | Default seats from `ENROLLMENT_CAP_DEFAULT` when body sends `0` |
| `.env.example` | Safe template (no real secrets) |
| `.gitignore` | Ignore `.env` |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged — still apply with `psql "$DATABASE_URL" -f db-ddl.sql`. Column remains **`seats`**.

## Env vars

| Variable | Required | Default |
|----------|----------|---------|
| `DATABASE_URL` | yes | — |
| `HTTP_ADDR` | no | `:8080` |
| `ENROLLMENT_CAP_DEFAULT` | no | `30` |
| `SHUTDOWN_TIMEOUT` | no | `10s` |

Do not log the full `DATABASE_URL` (password lives in the URL).

## Layout

```text
.env.example
db-ddl.sql
internal/config/config.go
cmd/api/main.go
```

## Run

```bash
git checkout lesson-25
go test ./internal/config ./internal/api ./internal/course
go build -o bin/api ./cmd/api

export DATABASE_URL='postgres://tayjava:tayjava@localhost:5432/tayjava?sslmode=disable'
export HTTP_ADDR=':8081'
psql "$DATABASE_URL" -f db-ddl.sql
./bin/api
```

Without `DATABASE_URL`, the process exits immediately with a clear error (expected).

HTTP tests and `go build ./...` do not need a live database.

## Try yourself

1. Unset `DATABASE_URL` and confirm fail-fast.
2. Set `ENROLLMENT_CAP_DEFAULT=50` and `POST` a course with `"seats":0` — response seats should be 50.
3. Copy `.env.example` → `.env` locally; keep `.env` untracked.

## What the next lesson adds

**Lesson 26** — multi-stage Dockerfile for a small Go API image.
