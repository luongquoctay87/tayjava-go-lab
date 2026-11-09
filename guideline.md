# Lesson 30 — Data races and go test -race

Branch: `lesson-30`  
Article: [Data races and go test -race in Golang](https://tayjava.com/vi/technologies/data-races-and-go-test-race-in-golang)

## Goal of this branch

Show a shared Course cache race and the mutex fix. Make `-race` part of the lab habit for concurrent packages.

`ReserveSeat` on the unsafe cache writes `Seats` with no lock. Concurrent map reads alone are not a data race; the race is that field write (and any map write such as `Put`). The safe cache uses `Seats`, not `capacity` or `seats_left`.

## Changes vs `lesson-29`

| Add / change | Why |
|--------------|-----|
| `internal/coursecache` | Unsafe vs Safe cache, `ReserveSeat`, concurrent test |
| `unsafe_race_test.go` | Build-tag demo that fails under `-race` |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged.

## Run

```bash
git checkout lesson-30
go test -race ./internal/coursecache/...
go test -race ./internal/notify/... ./internal/export/...
go build ./...
```

Deliberate race (expect FAIL):

```bash
go test -race -tags=racedemo ./internal/coursecache -run TestUnsafeCacheRace -count=1
```

## What the next lesson adds

**Lesson 31** — generics helpers (`Unique` / `First`) where they pay off.
