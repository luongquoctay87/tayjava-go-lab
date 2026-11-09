# Lesson 19 — HTTP server with net/http (handler + routing)

Branch: `lesson-19`  
Article: [HTTP server với net/http trong Golang — handler và routing](https://tayjava.com/vi/technologies/http-server-voi-nethttp-trong-golang-handler-va-routing)

## Goal of this branch

Run a Course API skeleton on Go 1.22+ `ServeMux` with method patterns and `PathValue`, returning plain text for health / course / enrollment.

## Changes vs `lesson-18`

| Add / change | Why |
|--------------|-----|
| `cmd/api/main.go` | Listen on `:8080`, wire mux |
| `internal/course/handler.go` | `Healthz`, `GetCourse`, `CreateEnrollment` |
| `README.md` | Branch table through lesson 19 |
| `guideline.md` | This file |

JSON file CLI store is unchanged. No Postgres / `db-ddl.sql` yet.

## Layout

```text
cmd/api/main.go
internal/course/
  model.go store.go …   # CLI
  handler.go            # HTTP skeleton ← new
```

## Run

```bash
git checkout lesson-19
go build -o bin/api ./cmd/api
./bin/api
```

In another terminal:

```bash
curl -i http://localhost:8080/healthz
curl -i http://localhost:8080/api/courses/101
curl -i "http://localhost:8080/api/enrollments?courseId=101&userId=7" -X POST
curl -i http://localhost:8080/api/enrollments   # expect 405
```

Expected (shape): `200 ok`, `course id=101`, `201 enrolled user=7 course=101`, method mismatch → 405.

Stop the server with Ctrl+C when done.

## Try yourself

1. Hit an unknown path — expect 404.
2. Log `r.Context()` cancellation by interrupting a slow client (prep for DB work).
3. Keep handlers thin — JSON + middleware arrive in lesson 20.

## What the next lesson adds

**Lesson 20** — JSON encode/decode plus logging/recovery middleware.
