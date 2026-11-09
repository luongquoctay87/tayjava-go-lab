# Lesson 20 — JSON and logging/recovery middleware

Branch: `lesson-20`  
Article: [JSON và middleware logging/recovery tự viết trong Golang](https://tayjava.com/vi/technologies/json-va-middleware-loggingrecovery-tu-viet-trong-golang)

## Goal of this branch

Encode/decode Course JSON, limit request bodies with `MaxBytesReader`, and wrap the mux with recovery (inner) then logging (outer) middleware written in stdlib.

## Changes vs `lesson-19`

| Add / change | Why |
|--------------|-----|
| `internal/api/middleware.go` | `Logging`, `Recovery`, `statusWriter` |
| `internal/api/middleware_test.go` | Panic → 500; logging status smoke test |
| `internal/course/handler.go` | JSON `CreateCourse` / `GetCourse` / enroll |
| `cmd/api/main.go` | Wire middleware stack + `POST /api/courses` |
| `README.md` | Branch table through lesson 20 |
| `guideline.md` | This file |

## Layout

```text
cmd/api/main.go
internal/api/middleware.go
internal/course/handler.go
```

## Run

```bash
git checkout lesson-20
go test ./internal/api ./internal/course
go build -o bin/api ./cmd/api
./bin/api
```

In another terminal (free port `:8080`):

```bash
curl -i -d '{"title":"Go Context","seats":30}' \
  -H 'Content-Type: application/json' \
  http://localhost:8080/api/courses
curl -i http://localhost:8080/api/courses/101
```

Watch the server log for `POST /api/courses 201 …`.

## Try yourself

1. Temporarily `panic("boom")` inside `CreateCourse` — client gets 500; process stays up.
2. Send a huge body and observe `MaxBytesReader` rejection.
3. Swap middleware order and compare logged status on panic.

## What the next lesson adds

**Lesson 21** — handler → service → repository (in-memory) ready for Postgres.
