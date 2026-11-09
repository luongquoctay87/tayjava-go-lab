# Lesson 15 — Channels: send and receive between goroutines

Branch: `lesson-15`  
Article: [Channel trong Golang — gửi nhận dữ liệu giữa goroutine](https://tayjava.com/vi/technologies/channel-trong-golang-gui-nhan-du-lieu-giua-goroutine)

## Goal of this branch

Share memory by communicating: unbuffered `chan EnrollmentEvent`, directional types (`chan<-` / `<-chan`), producer `close`, consumer `range`, and a `done` signal instead of `Sleep`.

## Changes vs `lesson-14`

| Add / change | Why |
|--------------|-----|
| `cmd/channel-lab/main.go` | Producer → notifier pipeline for enrollment events |
| `README.md` | Branch table through lesson 15 |
| `guideline.md` | This file |

## Layout

```text
cmd/
  goroutine-lab/   # lesson 14
  channel-lab/     # lesson 15 ← new
```

## Run

```bash
git checkout lesson-15
go run ./cmd/channel-lab
go build ./...
```

Expected:

```text
notify course=101 email=an@tayjava.com
notify course=101 email=binh@tayjava.com
notify course=102 email=chi@tayjava.com
all done
```

## Try yourself

1. Remove `close(out)` — `range` hangs forever (leak / block).
2. Drop `<-done` and see `main` exit before all notifies finish.
3. Sketch who owns the channel close in your design (lesson 16 digs into close rules).

## What the next lesson adds

**Lesson 16** — `select`, buffered channels, and close conventions (`cmd/select-lab`).
