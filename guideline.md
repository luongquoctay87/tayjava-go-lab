# Lesson 14 — Goroutines: run work concurrently the right way

Branch: `lesson-14`  
Article: [Goroutine trong Golang là gì? Chạy hàm đồng thời đúng cách](https://tayjava.com/vi/technologies/goroutine-trong-golang-la-gi-chay-ham-dong-thoi-dung-cach)

## Goal of this branch

Start goroutines with `go`, wait with `sync.WaitGroup`, and recognize why fire-and-forget (no join) drops work when `main` returns.

## Changes vs `lesson-13`

| Add / change | Why |
|--------------|-----|
| `cmd/goroutine-lab/main.go` | Sequential vs concurrent enrollment notify demo |
| `README.md` | Branch table through lesson 14 |
| `guideline.md` | This file |

Course CLI (`cmd/course-cli`, `internal/course`) stays unchanged.

Not yet: channels (15), `select` (16), Mutex deep dive (17), `context` (18), HTTP. No `db-ddl.sql` until lesson 22.

## Layout

```text
cmd/
  course-cli/      # lesson 13
  goroutine-lab/   # lesson 14 ← new
internal/course/   # CLI JSON store
```

## Run

```bash
git checkout lesson-14
go run ./cmd/goroutine-lab
go build ./...
```

Expected (shape — emails may print in any order in the concurrent section):

```text
--- sequential ---
enrolled course=101 email=an@tayjava.com
...
sequential done in ~250ms
--- concurrent (WaitGroup) ---
enrolled course=101 email=...
concurrent done in ~50ms
all notified
```

## Try yourself

1. Comment out `wg.Wait()` and re-run — `main` exits early; some lines vanish.
2. Compare sequential vs concurrent wall time for five 50ms sleeps.
3. Optionally add a tiny race demo package and run `go test -race` (full Mutex story in lesson 17).

## What the next lesson adds

**Lesson 15** — channels: send/receive `EnrollmentEvent` between producer and notifier goroutines.
