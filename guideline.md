# Lesson 35 — Retry, backoff, and idempotency

Branch: `lesson-35`  
Article: [Retry, backoff, and idempotency in Golang](https://tayjava.com/vi/technologies/retry-backoff-and-idempotency-in-golang)

## Goal of this branch

Add `Retry` with exponential backoff + jitter, classify temporary vs sold-out errors, and honor `Idempotency-Key` on enroll so retries do not double-enroll.

## Changes vs `lesson-34`

| Add / change | Why |
|--------------|-----|
| `internal/resilience` | `Retry`, `IdempotencyStore`, tests |
| `internal/api/enrollment.go` | Read `Idempotency-Key`; return prior enrollment |
| `db-ddl.sql` | Optional `idempotency_keys` table (durable path) |
| `cmd/api/main.go` | Wire in-memory store |
| `README.md` / `guideline.md` | This lesson |

### Schema note

`db-ddl.sql` now includes optional `idempotency_keys`. Existing Compose volumes that already initialized Postgres need a one-shot apply or volume reset to pick up the new table. The API lab path uses the in-memory store regardless.

`seats` remains total capacity.

## Run

```bash
git checkout lesson-35
go test ./internal/resilience/...
go test ./internal/api/...
go build ./...
```

## What the next lesson adds

**Lesson 36** — JWT Bearer middleware; secret from env (`JWT_SECRET`).
