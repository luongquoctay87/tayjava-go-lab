# Lesson 33 — Structured logging with slog

Branch: `lesson-33`  
Article: [Structured logging with slog in Golang](https://tayjava.com/vi/technologies/structured-logging-with-slog-in-golang)

## Goal of this branch

Use `log/slog` JSON as the default handler. Access logs carry `request_id`; enroll logs always include `course_id` / `user_id`. Never log JWTs or passwords.

## Changes vs `lesson-32`

| Add / change | Why |
|--------------|-----|
| `cmd/api/main.go` | `slog.SetDefault` JSON Info handler |
| `internal/api/middleware.go` | Structured access log + `X-Request-ID` |
| `internal/api/enrollment.go` | enroll_attempt / enroll_failed / list count |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged.

## Run

```bash
git checkout lesson-33
go test ./internal/api/...
go build -o bin/api ./cmd/api
```

## What the next lesson adds

**Lesson 34** — basic metrics counters / OpenTelemetry concepts (light demo).
