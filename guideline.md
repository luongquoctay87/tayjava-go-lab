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

- Workers honor cancel or exit when jobs channel closes.
- Only one place closes the jobs channel (`Close`).
- Do not share unsynchronized maps across workers — run `go test -race ./internal/notify/...`.
- Notify uses `context.Background()` so client cancel does not drop accepted work.

## What the next lesson adds

**Lesson 29** — `errgroup`, pipelines, and graceful HTTP shutdown (`Server.Shutdown` then `pool.Close`).
