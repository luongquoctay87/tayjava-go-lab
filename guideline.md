# Lesson 22 — Postgres with database/sql + db-ddl.sql

Branch: `lesson-22`  
Article: [Postgres với database/sql trong Golang — và migration cơ bản](https://tayjava.com/vi/technologies/postgres-voi-databasesql-trong-golang-va-migration-co-ban)

## Goal of this branch

Open a `*sql.DB` pool with `PingContext`, implement `PostgresRepo` for Course CRUD, and introduce the shared Postgres schema file `db-ddl.sql`. Capacity column is **`seats`** (not `seats_left`).

## Changes vs `lesson-21`

| Add / change | Why |
|--------------|-----|
| `db-ddl.sql` | `courses` + `enrollments` schema (kept from this lesson onward) |
| `internal/course/postgres.go` | `OpenDB`, `PostgresRepo.Get` / `Save` |
| `cmd/api/main.go` | Wire Postgres when `DATABASE_URL` is set |
| `go.mod` / `go.sum` | `github.com/lib/pq` driver |
| `README.md` | Branch table through lesson 22 |
| `guideline.md` | This file |

Service and HTTP handlers stay the same — only the repository implementation changes.

## Schema notes

- `courses.seats` = available capacity (`CHECK (seats >= 0)`).
- Optional forward-compat columns: `level`, `created_at` (CLI / later labs).
- `enrollments`: FK to `courses`, `UNIQUE (course_id, user_id)`.

Apply DDL:

```bash
psql "$DATABASE_URL" -f db-ddl.sql
```

Example DSN:

```text
postgres://tayjava:tayjava@localhost:5432/tayjava?sslmode=disable
```

## Layout

```text
db-ddl.sql
cmd/api/main.go
internal/course/postgres.go   # OpenDB + PostgresRepo
internal/course/repository.go # MemoryRepo still available
```

## Run

```bash
git checkout lesson-22
go test ./internal/course ./internal/api
go build -o bin/api ./cmd/api
```

Without Postgres (memory fallback):

```bash
./bin/api
```

With Postgres:

```bash
export DATABASE_URL='postgres://tayjava:tayjava@localhost:5432/tayjava?sslmode=disable'
psql "$DATABASE_URL" -f db-ddl.sql
./bin/api
```

```bash
curl -i -d '{"title":"Go Context","seats":30,"level":"intermediate"}' \
  -H 'Content-Type: application/json' \
  http://localhost:8080/api/courses
curl -i http://localhost:8080/api/courses/1
```

Restart the API — rows remain when using Postgres (unlike MemoryRepo).

## Try yourself

1. Confirm `$1` placeholders in `PostgresRepo` — never concatenate user input into SQL.
2. Drop `courses` and restart without re-applying DDL — expect a clear SQL error.
3. Swap back to `NewMemoryRepo()` in `main` and note data disappears after restart.

## What the next lesson adds

**Lesson 23** — validation at the HTTP edge and centralized `MapDomainError` JSON errors.
