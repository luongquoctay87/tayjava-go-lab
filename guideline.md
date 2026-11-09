# Lesson 21 — Handler → service → repository

Branch: `lesson-21`  
Article: [Tổ chức handler → service → repository trong Golang](https://tayjava.com/vi/technologies/to-chuc-handler-service-repository-trong-golang)

## Goal of this branch

Split the Course API into thin HTTP handlers, a use-case `Service`, and a small `Repository` interface with an in-memory implementation — ready to swap for Postgres in lesson 22. No `db-ddl.sql` yet.

## Changes vs `lesson-20`

| Add / change | Why |
|--------------|-----|
| `internal/course/repository.go` | `Repository` + `MemoryRepo` |
| `internal/course/service.go` | `Enroll` / `GetCourse` / `CreateCourse` + domain errors |
| `internal/course/service_test.go` | Table-driven enroll tests (no HTTP) |
| `internal/api/course.go` | Course HTTP handlers |
| `internal/api/enrollment.go` | Enrollment HTTP handlers |
| `internal/api/handler_test.go` | End-to-end HTTP against memory repo |
| `cmd/api/main.go` | Explicit constructor wiring |
| removed `internal/course/handler.go` | HTTP lives under `internal/api` |
| `README.md` | Branch table through lesson 21 |
| `guideline.md` | This file |

JSON file CLI (`cmd/course-cli` + `Store`) stays for the beginner lab path.

## Layout

```text
cmd/api/main.go
internal/course/
  model.go store.go          # CLI JSON
  repository.go service.go   # API domain layers
internal/api/
  middleware.go
  course.go enrollment.go
```

## Run

```bash
git checkout lesson-21
go test ./internal/course ./internal/api
go build -o bin/api ./cmd/api
./bin/api
```

Seeded memory course: `id=1`, `seats=2`.

```bash
curl -i http://localhost:8080/api/courses/1
curl -i -d '{"courseId":1,"userId":7}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
curl -i -d '{"courseId":1,"userId":8}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
curl -i -d '{"courseId":1,"userId":9}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments   # expect 409 no seats left
```

## Try yourself

1. Enroll the same `userId` twice — expect `409 already enrolled`.
2. Implement a second `Repository` (e.g. map copy) and swap it in `main` — service unchanged.
3. Keep each handler under ~40 lines; map domain errors only at the HTTP edge.

## What the next lesson adds

**Lesson 22** — Postgres with `database/sql` and the first `db-ddl.sql` schema.
