# Lesson 28 — Worker pools and fan-in / fan-out

Branch: `lesson-28`  
Article: [Worker pools and fan-in / fan-out in Golang](https://tayjava.com/vi/technologies/worker-pools-and-fan-in-fan-out-in-golang)

## Goal of this branch

Add a bounded notify worker pool after successful enrollment. Cap goroutines (workers=4, queue=64); demonstrate fan-in aggregation; keep `db-ddl.sql` unchanged.

## Changes vs `lesson-27`

| Add / change | Why |
|--------------|-----|
| `internal/notify` | Fixed-size pool + `FanIn` helper + tests |
| `internal/api/enrollment.go` | Submit notify job after successful Enroll |
| `cmd/api/main.go` | Wire pool; `defer pool.Close()` |
| `README.md` / `guideline.md` | This lesson |

## Seats / schema

`db-ddl.sql` is unchanged. Column **`seats`** remains total capacity.

## Run

```bash
git checkout lesson-28
go test ./internal/notify/...
go test ./...
go build -o bin/api ./cmd/api
```

With Compose (from lesson 27):

```bash
docker compose up --build
# enroll — watch fake SMTP logs: enrolled user=… course=…
```

## Leak / race checklist

- `Submit` returns an error: `ErrClosed` after `Close`, or `ctx.Err()` if the queue wait times out. It does not panic on a closed channel.
- `Close` stops new jobs and drains the queue. Canceling the pool context drops jobs still queued.
- The enroll handler waits at most 100ms for a queue slot. Once queued, the pool context owns the job.
- Only one place closes the jobs channel (`Close`, via `sync.Once`).
- Do not share unsynchronized maps across workers — run `go test -race ./internal/notify/...`.

## What the next lesson adds

**Lesson 29** — `errgroup`, pipelines, and graceful HTTP shutdown (`Server.Shutdown` then `pool.Close`).
