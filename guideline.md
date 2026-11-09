# Lesson 18 — context: timeout, cancel, and call-chain propagation

Branch: `lesson-18`  
Article: [context.Context trong Golang — timeout, cancel và truyền xuống call chain](https://tayjava.com/vi/technologies/contextcontext-trong-golang-timeout-cancel-va-truyen-xuong-call-chain)

## Goal of this branch

Create contexts with `WithTimeout` / `WithCancel`, always `defer cancel()`, pass `ctx` through Enroll → checkSeat → saveEnrollment, and carry a request ID via `WithValue` (not a DB handle).

## Changes vs `lesson-17`

| Add / change | Why |
|--------------|-----|
| `cmd/contextlab/main.go` | Timeout, cancel, and Enroll chain demos |
| `README.md` | Branch table through lesson 18 |
| `guideline.md` | This file |

## Layout

```text
cmd/contextlab/main.go
```

## Run

```bash
git checkout lesson-18
go run ./cmd/contextlab
go build ./...
```

Expected (shape):

```text
--- timeout too short (expect deadline) ---
failed: context deadline exceeded
--- enroll with request ID (expect success) ---
enrolled user=7 course=101 requestID=req-42
enroll ok
--- enroll deadline (80ms+80ms > 100ms) ---
enroll failed: context deadline exceeded
--- manual cancel ---
result:  context canceled
```

## Try yourself

1. Raise the first timeout above 800ms — `loadCourse` succeeds.
2. Omit `defer cancel()` once while reading docs — then put it back (timer leak habit).
3. Never create `context.Background()` inside `saveEnrollment`.

## What the next lesson adds

**Lesson 19** — `net/http` Course API skeleton (`cmd/api`).
