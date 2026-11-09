# Lesson 22 — Postgres with database/sql and golang-migrate

Branch: `lesson-22`  
Article: [Postgres with database/sql in Golang — and basic migrations](https://tayjava.com/vi/technologies/postgres-voi-databasesql-trong-golang-va-migration-co-ban)

## Goal of this branch

Open a `*sql.DB` pool with the `pgx` stdlib driver, `PingContext`, and `SetMaxIdleConns`. `PostgresRepo` loads and saves courses and enrolls with a conditional `UPDATE`. Schema changes live in golang-migrate `.up.sql` / `.down.sql` files. `db-ddl.sql` stays in sync with those files for a one-shot `psql` apply. Capacity column is **`seats`** (not `seats_left`).

`Service.Enroll` is unchanged from lesson 21: it checks that the course exists, then calls `repo.Enroll`. Handlers stay the same. `main` reads `DATABASE_URL` and uses `MemoryRepo` when that variable is empty.

## Changes vs `lesson-21`

| Add / change | Why |
|--------------|-----|
| `migrations/000001_create_courses.*.sql` | Versioned `courses` table |
| `migrations/000002_create_enrollments.*.sql` | Versioned `enrollments` table and unique pair |
| `db-ddl.sql` | Same tables, for `psql -f` without the migrate CLI |
| `internal/course/postgres.go` | `OpenDB`, `PostgresRepo.Get` / `Save` / `Enroll` |
| `cmd/api/main.go` | `DATABASE_URL` selects Postgres vs memory |
| `go.mod` / `go.sum` | `github.com/jackc/pgx/v5` v5.7.4 (works with Go 1.22) |
| `README.md` | Branch table through lesson 22 |
| `guideline.md` | This file |

## How enroll works

`PostgresRepo.Enroll` does **not** read `seats` and write it back. It runs:

```sql
UPDATE courses SET seats = seats - 1 WHERE id = $1 AND seats > 0
```

`RowsAffected() == 0` means no row had seats left, so the repo returns `ErrNoSeats`. The service already returned `ErrNotFound` when the course id does not exist.

Then it inserts into `enrollments`. Postgres unique violation `23505` (`UNIQUE (course_id, user_id)`) becomes `ErrConflict`.

These two statements are **not** in a transaction. If the insert fails, the seat is already decremented and there is no matching enrollment. Lesson 38 wraps them in `BeginTx` / `Commit`.

## Pool

`OpenDB` uses `sql.Open("pgx", dsn)`. The blank import of `github.com/jackc/pgx/v5/stdlib` is in `postgres.go`, next to `OpenDB`. `SetMaxOpenConns(10)` and `SetMaxIdleConns(10)` keep idle connections from dropping to the `database/sql` default of 2. `PingContext` runs before the pool is returned. The DSN is `DATABASE_URL` only — no password in source.

Newer pgx v5 releases need a Go newer than 1.22. This module stays on `go 1.22` and pins **v5.7.4**.

## Schema notes

- `courses.seats` is available capacity (`INT`, `CHECK (seats >= 0)`).
- `level` and `created_at` are on `courses` for the CLI and later labs.
- `enrollments.course_id` / `user_id` are `BIGINT`. Unique on `(course_id, user_id)`.
- IDs are `BIGSERIAL` (`int64` in Go).

Apply with golang-migrate (preferred):

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
export DATABASE_URL='postgres://tayjava:<local-password>@localhost:5432/courses?sslmode=disable'
migrate -path migrations -database "$DATABASE_URL" up
```

Or apply the same SQL in one step:

```bash
psql "$DATABASE_URL" -f db-ddl.sql
```

`migrate ... down 1` runs the latest `.down.sql`. Do not edit the two files out of sync.

## Layout

```text
db-ddl.sql
migrations/
  000001_create_courses.up.sql
  000001_create_courses.down.sql
  000002_create_enrollments.up.sql
  000002_create_enrollments.down.sql
cmd/api/main.go
internal/course/postgres.go   # OpenDB + PostgresRepo
internal/course/repository.go # MemoryRepo still available for tests
```

## Run

```bash
git checkout lesson-22
go test ./internal/course ./internal/api
go build -o bin/api ./cmd/api
```

Tests use `MemoryRepo` and do not need a database. `go build` still compiles `PostgresRepo`.

Without Postgres (memory fallback):

```bash
./bin/api
```

With Postgres:

```bash
export DATABASE_URL='postgres://tayjava:<local-password>@localhost:5432/courses?sslmode=disable'
migrate -path migrations -database "$DATABASE_URL" up
./bin/api
```

```bash
curl -i -d '{"title":"Go Context","seats":1,"level":"intermediate"}' \
  -H 'Content-Type: application/json' \
  http://localhost:8080/api/courses
curl -i http://localhost:8080/api/courses/1
```

Restart the API — rows remain when using Postgres (unlike MemoryRepo).

Seed a course with `seats = 1` and post several enrollments with different `userId` values. Only one succeeds, and `seats` does not go negative. Enrolling the same user again returns `ErrConflict` (the seat from the failed insert is not restored until lesson 38).

## Try yourself

1. Confirm `$1` placeholders in `PostgresRepo` — never concatenate user input into SQL.
2. Drop `courses` and restart without migrating — expect a clear SQL error.
3. Unset `DATABASE_URL` and note data disappears after restart (`MemoryRepo`).

## What the next lesson adds

**Lesson 23** — validation at the HTTP edge and centralized `MapDomainError` JSON errors.
