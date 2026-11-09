# Lesson 29 — errgroup, pipelines, and graceful shutdown

Branch: `lesson-29`  
Article: [errgroup, pipelines, and graceful shutdown in Golang](https://tayjava.com/vi/technologies/errgroup-pipelines-and-graceful-shutdown-in-golang)

## Goal of this branch

Wire `errgroup` for a 3-stage export pipeline and graceful HTTP shutdown: signal → `Server.Shutdown` → `pool.Close()` drain. Pool still uses `context.Background()` so workers do not exit before the queue drains.

## Changes vs `lesson-28`

| Add / change | Why |
|--------------|-----|
| `internal/export` | errgroup pipeline demo + tests |
| `cmd/api/main.go` | `signal.NotifyContext` + `Shutdown` + `pool.Close` |
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

Shutdown smoke: start the API, enroll a few users, then Ctrl+C — expect “shutting down…” then “bye” without a hung process.

## What the next lesson adds

**Lesson 30** — data races and `go test -race` on a shared cache demo.
