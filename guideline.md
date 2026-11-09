# Lesson 29 — errgroup, pipelines, and graceful shutdown

Branch: `lesson-29`  
Article: [errgroup, pipelines, and graceful shutdown in Golang](https://tayjava.com/vi/technologies/errgroup-pipelines-and-graceful-shutdown-in-golang)

## Goal of this branch

Wire `errgroup` for a 3-stage export pipeline and graceful HTTP shutdown: signal → `stop()` so a second signal exits → `Server.Shutdown` with a deadline → on success, `pool.Close()` bounded by that same deadline. If `Shutdown` returns an error, skip the drain. The pool uses `context.Background()` so workers are not canceled by the signal before `Close` can drain.

The fetch stage uses `errgroup.WithContext` so the first failure cancels sibling fetches. Go 1.22 gives each loop iteration its own variable, so the lab does not rebind `id := id`.

## Changes vs `lesson-28`

| Add / change | Why |
|--------------|-----|
| `internal/export` | errgroup pipeline demo + tests |
| `cmd/api/main.go` | `signal.NotifyContext`, `stop()` on the first signal, `Shutdown`, then a bounded `pool.Close` |
| `go.mod` | `golang.org/x/sync` |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged (`seats` = total capacity).

## Run

```bash
git checkout lesson-29
go test ./internal/export/...
go test ./...
go build -o bin/api ./cmd/api
```

Shutdown smoke: start the API, enroll a few users, then Ctrl+C — expect “shutting down…” then “bye”. A second Ctrl+C exits immediately. If HTTP shutdown hits the deadline, the process logs the error and skips the pool drain.

## What the next lesson adds

**Lesson 30** — data races and `go test -race` on a shared cache demo.
