# Lesson 24 — HTTP tests with httptest

Branch: `lesson-24`  
Article: [Test HTTP trong Golang bằng httptest](https://tayjava.com/vi/technologies/test-http-trong-golang-bang-httptest)

## Goal of this branch

Lock the lesson-23 error contract with table-driven `httptest` tests, a stub `Repository` that implements `Enroll`, a mux routing smoke test, and a canceled-context case — no real listen port and no Postgres required. The unknown-field case sends `"seat"` and expects `400 VALIDATION` because `decodeEnroll` uses `DisallowUnknownFields`. The canceled-context test builds the request with `httptest.NewRequest` and `req.WithContext`, then expects `500 INTERNAL` (there is no `context.Canceled` branch in `MapDomainError`).

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
go test -race ./internal/api/ -v
go test -race ./...
go build -o bin/api ./cmd/api
```

## Try yourself

1. Flip expected status in a subtest and watch the failure message (status + body).
2. Remove `DisallowUnknownFields` and watch the unknown-field case stop returning `VALIDATION`.
3. Call `t.Parallel()` carefully — use a fresh stubRepo per subtest.

## What the next lesson adds

**Lesson 25** — env-based `Config` (twelve-factor light): `DATABASE_URL`, `HTTP_ADDR`.
