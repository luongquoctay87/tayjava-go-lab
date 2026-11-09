# Lesson 38 — Database transactions and context deadlines

Branch: `lesson-38`  
Article: [Database transactions and context deadlines in Golang](https://tayjava.com/vi/technologies/database-transactions-and-context-deadlines-in-golang)

## Goal of this branch

Enroll inside one transaction: `SELECT seats … FOR UPDATE`, count enrollments, insert when `count < seats`. Handler uses `context.WithTimeout(3s)`. Column remains **`seats`** (total capacity — not `seats_left`).

## Changes vs `lesson-37`

| Add / change | Why |
|--------------|-----|
| `PostgresRepo.EnrollTx` | BeginTx + FOR UPDATE + capacity check |
| `MemoryRepo.EnrollTx` | Mutex stand-in for unit tests |
| `Service.Enroll` | Delegates to `EnrollTx` |
| `EnrollmentHandler` | 3s timeout → 504 on deadline |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged (still includes optional `idempotency_keys` from lesson 35).

## Run

```bash
git checkout lesson-38
go test ./internal/course/... ./internal/api/...
go build ./...
```

## What the next lesson adds

**Lesson 39** — Advanced lab: API accept → bounded enrollment worker → metrics + graceful drain.
