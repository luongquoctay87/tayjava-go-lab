# Lesson 24 — HTTP tests with httptest

Branch: `lesson-24`  
Article: [Test HTTP trong Golang bằng httptest](https://tayjava.com/vi/technologies/test-http-trong-golang-bang-httptest)

## Goal of this branch

Lock the lesson-23 error contract with table-driven `httptest` tests, a stub `Repository`, mux routing smoke test, and a canceled-context case — no real listen port and no Postgres required.

## Changes vs `lesson-23`

| Add / change | Why |
|--------------|-----|
| `internal/api/enrollment_test.go` | Table-driven enroll tests + stubRepo + mux/context cases |
| `README.md` | Branch table through lesson 24 |
| `guideline.md` | This file |

Handlers and `db-ddl.sql` unchanged. Unit tests stay on MemoryRepo / stubRepo.

## Layout

```text
db-ddl.sql
internal/api/enrollment_test.go
internal/api/handler_test.go
internal/api/errors.go
```

## Run

```bash
git checkout lesson-24
go test ./internal/api/ -v
go test ./...
go build -o bin/api ./cmd/api
```

## Try yourself

1. Flip expected status in a subtest and watch the failure message (status + body).
2. Add a case for unknown JSON fields (`DisallowUnknownFields`).
3. Call `t.Parallel()` carefully — use a fresh stubRepo per subtest.

## What the next lesson adds

**Lesson 25** — env-based `Config` (twelve-factor light): `DATABASE_URL`, `HTTP_ADDR`.
