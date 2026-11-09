# Lesson 17 — WaitGroup and Mutex when channels are not enough

Branch: `lesson-17`  
Article: [sync.WaitGroup và Mutex trong Golang — khi channel chưa đủ](https://tayjava.com/vi/technologies/syncwaitgroup-va-mutex-trong-golang-khi-channel-chua-du)

## Goal of this branch

Join N workers with `WaitGroup`, protect shared seat counts with `Mutex`, and keep `go test -race` clean on concurrent reserves.

## Changes vs `lesson-16`

| Add / change | Why |
|--------------|-----|
| `internal/seatledger/` | Mutex-backed `SeatLedger` + race-friendly tests |
| `cmd/synclab/main.go` | WaitGroup notify + concurrent reserve demo |
| `README.md` | Branch table through lesson 17 |
| `guideline.md` | This file |

Note: the lesson-13 Course CLI is single-threaded — it does not need a Mutex yet.

## Layout

```text
internal/seatledger/
  ledger.go
  ledger_test.go
cmd/synclab/
```

## Run

```bash
git checkout lesson-17
go test -race ./internal/seatledger
go run ./cmd/synclab
go build ./...
```

Expected test: pass with race detector clean.

Expected `synclab` (reserve lines may reorder):

```text
--- WaitGroup notify ---
sent: ...
all emails sent
--- Mutex seat ledger ---
reserved true|false (exactly two true)
left 0
```

## Try yourself

1. Temporarily comment `Lock`/`Unlock` in `Reserve`, run `go test -race` — expect a WARNING, then restore.
2. Forget one `Done` and watch `Wait` hang.
3. Confirm CLI Course CRUD still needs no Mutex (one goroutine).

## What the next lesson adds

**Lesson 18** — `context` timeout/cancel through a fake Enroll call chain.
