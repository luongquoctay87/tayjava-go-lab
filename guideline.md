# Lesson 27 — Intermediate lab: Course REST + enrollment

Branch: `lesson-27`  
Article: [Lab Intermediate: Course REST API + enrollment bằng Golang](https://tayjava.com/vi/technologies/lab-intermediate-course-rest-api-enrollment-bang-golang)

## Goal of this branch

Compose lessons 19–26 into an end-to-end Intermediate lab: Course REST + enrollment capped by **`seats`** (total capacity, not `seats_left`), Postgres persistence, env config, multi-stage Docker image, and `docker compose` with schema init from `db-ddl.sql`.

## Changes vs `lesson-26`

| Add / change | Why |
|--------------|-----|
| `docker-compose.yml` | API + Postgres; DDL via init volume |
| `internal/course/*` | Enrollment persistence; seats = capacity |
| `internal/api/enrollment.go` | Nested `/courses/{id}/enrollments` |
| `cmd/api/main.go` | Intermediate routes + DB retry for Compose |
| `README.md` / `guideline.md` | This lesson |

## Seats semantics

Column / JSON field remains **`seats`**. From this lesson it means **total capacity**. Enrollment rows are counted separately; when `count >= seats` the API returns `409 NO_SEATS`. Do not rename to `seats_left` / `capacity` in this lab.

Apply DDL (Compose does this on first `db` boot):

```bash
psql "$DATABASE_URL" -f db-ddl.sql
```

## Endpoints

| Method | Path | Notes |
|--------|------|-------|
| `GET` | `/healthz` | Liveness |
| `POST` | `/courses` | Create (`title`, `seats`, optional `level`) |
| `GET` | `/courses/{id}` | Detail |
| `POST` | `/courses/{id}/enrollments` | Body `{"userId":…}` |
| `GET` | `/courses/{id}/enrollments` | List |
| `POST` | `/api/courses`, `/api/enrollments` | Aliases from earlier lessons |

## Layout

```text
Dockerfile
docker-compose.yml
db-ddl.sql
.env.example
cmd/api/main.go
internal/config/
internal/course/
internal/api/
```

## Run

Unit tests (no Docker / DB required):

```bash
git checkout lesson-27
go test ./...
go build -o bin/api ./cmd/api
```

Compose end-to-end:

```bash
docker compose up --build
```

Sample curls (another terminal):

```bash
curl -s -d '{"title":"Go Intermediate","seats":2}' -H 'Content-Type: application/json' \
  http://localhost:8080/courses
curl -s -d '{"userId":7}' -H 'Content-Type: application/json' \
  http://localhost:8080/courses/1/enrollments
curl -s -d '{"userId":8}' -H 'Content-Type: application/json' \
  http://localhost:8080/courses/1/enrollments
curl -s -d '{"userId":9}' -H 'Content-Type: application/json' \
  http://localhost:8080/courses/1/enrollments
# expect 409 NO_SEATS
curl -s http://localhost:8080/courses/1/enrollments
```

Without Compose, export `DATABASE_URL`, apply `db-ddl.sql`, then `go run ./cmd/api`.

## Try yourself

1. Enroll the same `userId` twice — expect `409 CONFLICT`.
2. Note concurrent enroll races (no `SELECT … FOR UPDATE` yet) — Advanced lessons dig in.
3. Confirm `.env` stays gitignored; only `.env.example` is committed.

## What the next lesson adds

**Lesson 28** — worker pool and fan-in / fan-out for parallel enroll/notify work.
