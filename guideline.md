# Lesson 39 — Advanced lab: API + enrollment worker + metrics

Branch: `lesson-39`  
Article: [Advanced Golang lab: API + enrollment worker + metrics](https://tayjava.com/vi/technologies/advanced-golang-lab-api-enrollment-worker-metrics)

## Goal of this branch

Wire one binary: HTTP validates + enqueues (`202`, `409` for a duplicate key, `503` with `Retry-After` when the queue is full) → bounded worker pool → seat write via lesson-38 transactions → backlog metrics on `GET /metrics` (`expvar`). Shutdown order: `srv.Shutdown` (check the error) → `Close` (stop enqueue, then close the job channel) → workers drain with `for range` → `Wait` bounded by the same shutdown deadline. Workers use their own context, not the signal context. A duplicate idempotency key is rejected before enqueue.

## Changes vs `lesson-38`

| Add / change | Why |
|--------------|-----|
| `internal/enrollment` | Enroller, async handler, metrics, tests |
| `cmd/api/main.go` | Async `POST /api/enrollments` + worker drain |
| `GET /metrics` | `expvar` counters and `enroll_queue_depth` |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged. Sync nested enroll (`POST /courses/{id}/enrollments`) still exists for Intermediate curls.

## Run

```bash
git checkout lesson-39
go test -race ./internal/enrollment/...
go test ./...
go build -o bin/api ./cmd/api
```

Async enroll (JWT required):

```bash
TOKEN=$(curl -s -d '{"userId":"7","password":"lab"}' -H 'Content-Type: application/json' localhost:8080/login | jq -r .token)
curl -s -o /dev/null -w '%{http_code}\n' \
  -d '{"courseId":1,"userId":7}' \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Idempotency-Key: abc' \
  -H 'Content-Type: application/json' \
  localhost:8080/api/enrollments
# 202
curl -s localhost:8080/metrics
```

## What the next lesson adds

**Lesson 40** — package design / `internal` boundaries / short compatibility ADR.
