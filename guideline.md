# Lesson 35 — Retry, backoff, and idempotency

Branch: `lesson-35`  
Article: [Retry, backoff, and idempotency in Golang](https://tayjava.com/vi/technologies/retry-backoff-and-idempotency-in-golang)

## Goal of this branch

Add `Retry` with exponential backoff and full jitter (`math/rand/v2`, capped by `maxDelay`). Classify temporary errors, including 429 and 503 that carry `Retry-After`. Honor `Idempotency-Key` per user: SHA-256 fingerprint of the raw body, replay of the original `201` and body, `422` when the same key is reused with a different body, and `409` with `Retry-After: 1` while the first call is in flight. The store mutex does not cover seat reservation.

## Changes vs `lesson-34`

| Add / change | Why |
|--------------|-----|
| `internal/resilience` | `Retry` (full jitter + cap), `Idempotency.Do`, tests |
| `internal/api/enrollment.go` | Per-user key, body fingerprint, replay `201` / `422` / `409` |
| `db-ddl.sql` | Optional `idempotency_keys` unique on `(user_id, key)` |
| `cmd/api/main.go` | Wire the in-memory store |
| `README.md` / `guideline.md` | This lesson |

When the handler is constructed with a store, `Idempotency-Key` is required (`400` if missing). The key is scoped as `userId:key`. `userId` still comes from the JSON body on this branch; lesson 36 moves it to the JWT. A `5xx` result is not stored, so a later retry with the same key may run again. Requests that pass `nil` for the store (older handler tests) enroll once without a key.

### Schema note

`db-ddl.sql` includes optional `idempotency_keys` with primary key `(user_id, key)`. Existing Compose volumes that already initialized Postgres need a one-shot apply or volume reset to pick up the new table. The API lab path uses the in-memory store regardless.

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
