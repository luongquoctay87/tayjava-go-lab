# Lesson 33 — Structured logging with slog

Branch: `lesson-33`  
Article: [Structured logging with slog in Golang](https://tayjava.com/vi/technologies/structured-logging-with-slog-in-golang)

## Goal of this branch

Use `log/slog` JSON as the default handler. Access logs carry `request_id`; enroll logs always include `course_id` / `user_id`. Never log JWTs or passwords.

`main` calls `slog.SetDefault` first, then passes that logger into `NewEnrollmentHandler`, which builds the `component=enrollment` child. Do not store `slog.Default().With(...)` in a package variable: init runs before `SetDefault`.

`X-Request-ID` is kept only when it is 1–64 characters of letters, digits, `-`, or `_`. Anything else is replaced with a generated id.

No seats left and a duplicate enrollment are `Info` (`reason=no_seats` or `conflict`). Unexpected failures, including a database error, are `Error`.

## Changes vs `lesson-32`

| Add / change | Why |
|--------------|-----|
| `cmd/api/main.go` | `slog.SetDefault`, then inject that logger into the enrollment handler |
| `internal/api/middleware.go` | Structured access log + validated `X-Request-ID` |
| `internal/api/enrollment.go` | Child logger; no-seats/conflict at Info, system failures at Error |
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
