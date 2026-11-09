# Lesson 27 — Intermediate lab: Course REST + enrollment

Branch: `lesson-27`  
Article: [Intermediate lab: Course REST API + enrollment in Golang](https://tayjava.com/en/technologies/intermediate-lab-course-rest-api-enrollment-in-golang)

## Goal of this branch

Assemble lessons 19–26 into an end-to-end Intermediate lab: create and get a Course, enroll with a seat limit, Postgres, env config (`HTTP_ADDR`, `DATABASE_URL`), the lesson-26 multi-stage image, and Docker Compose.

## Changes vs `lesson-26`

| Add / change | Why |
|--------------|-----|
| `docker-compose.yml` | Postgres healthcheck, migrate, then API |
| `migrations/` | Lesson 22 golang-migrate up/down for `courses` and `enrollments` |
| `internal/course` | `Enroll` uses `UPDATE … SET seats = seats - 1 WHERE seats > 0` (no transaction yet) |
| `cmd/api` | Canonical routes only; `-healthcheck` from lesson 26 |
| `.env.example` | Dummy `POSTGRES_PASSWORD`; real `.env` stays gitignored |

Column name remains **`seats`** (remaining capacity), not `seats_left`. Two statements run without a transaction: if the INSERT fails after the UPDATE, the seat is already gone. Lesson 38 wraps both in one transaction.

Only `ErrNotFound` from `Get` stays a 404. Any other `Get` error is wrapped with `%w` and `MapDomainError` returns 500.

## Endpoints

| Method | Path | Notes |
|--------|------|-------|
| `GET` | `/healthz` | Liveness |
| `POST` | `/api/courses` | Create (`title`, `seats`) |
| `GET` | `/api/courses/{id}` | Detail |
| `POST` | `/api/enrollments` | Body `{"courseId":1,"userId":7}` |

## Layout

```text
Dockerfile
docker-compose.yml
migrations/
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
go test -race ./...
go build -o bin/api ./cmd/api
```

Compose (copy `.env.example` to `.env` and set `POSTGRES_PASSWORD`):

```bash
docker compose up --build
```

Sample curls:

```bash
curl -s -d '{"title":"Go Intermediate","seats":2}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/courses
curl -s -d '{"courseId":1,"userId":7}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
curl -s -d '{"courseId":1,"userId":8}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
curl -s -d '{"courseId":1,"userId":9}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
# expect 409 NO_SEATS
curl -s http://localhost:8080/api/courses/1
```

Postgres is healthy before migrate (`pg_isready`). The API starts only after migrate exits successfully. The API healthcheck runs `["/app", "-healthcheck"]`. If the password contains `@` or `/`, URL-encode it in `DATABASE_URL`.

Without Compose, export `DATABASE_URL`, apply `db-ddl.sql` or `migrate -path migrations -database "$DATABASE_URL" up`, then `go run ./cmd/api`.

## Try yourself

1. Enroll the same `userId` twice — expect `409 CONFLICT`.
2. Stop Postgres and call enroll — expect `500 INTERNAL`, not `404 NOT_FOUND`.
3. Confirm `.env` stays gitignored; only `.env.example` is committed.

## What the next lesson adds

**Lesson 28** — worker pool and fan-in / fan-out for parallel enroll/notify work.
