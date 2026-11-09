# Lesson 16 — select, buffered channels, and close

Branch: `lesson-16`  
Article: [select, buffered channel và đóng channel trong Golang](https://tayjava.com/vi/technologies/select-buffered-channel-va-dong-channel-trong-golang)

## Goal of this branch

Use a small buffer for burst sends, `select` among ready cases (including `time.After` timeout), and close channels only from the send owner.

## Changes vs `lesson-15`

| Add / change | Why |
|--------------|-----|
| `cmd/select-lab/main.go` | Buffered jobs, select timeout, drain-on-close |
| `README.md` | Branch table through lesson 16 |
| `guideline.md` | This file |

## Layout

```text
cmd/
  channel-lab/   # lesson 15
  select-lab/    # lesson 16 ← new
```

## Run

```bash
git checkout lesson-16
go run ./cmd/select-lab
go build ./...
```

Expected (shape):

```text
--- drain buffered jobs ---
process course 101
process course 102
process course 103
--- select enroll / cancel / timeout ---
enrolled an@tayjava.com
--- select until closed ---
got course-101
got course-102
drained
```

## Try yourself

1. Raise buffer to 3 and send three IDs from `main` before starting a receiver.
2. Shorten enroll sleep so timeout still loses — or lengthen it to hit `timeout waiting enroll`.
3. Never double-`close` in real code (panic). Note who owns close: the producer.

## What the next lesson adds

**Lesson 17** — `sync.WaitGroup` and `Mutex` with a seat ledger (`internal/seatledger`).
