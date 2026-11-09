# Lesson 11 — Pointers when you need them

Branch: `lesson-11`  
Article: [Pointer đủ dùng trong Golang](https://tayjava.com/vi/technologies/pointer-du-dung-trong-golang)

## Goal of this branch

Mutate `Course.Seats` through a pointer (`Reserve`), contrast with pass-by-value (`ReserveCopy`), and handle `nil` safely.

## Changes vs `lesson-10`

| Add / change | Why |
|--------------|-----|
| `course/reserve.go` | `Reserve(*Course)` + `ReserveCopy(Course)` |
| `cmd/ptrdemo/main.go` | Runnable pointer demo |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-11
go run ./cmd/ptrdemo
go test ./...
```

Expected (shape):

```text
before reserve: 30
after reserve: 29
after reserveCopy (unchanged): 29 (was 29 )
reserve nil: course is nil
```

## Try yourself

1. Call `Reserve` until seats hit 0 and observe the error.
2. Pass `c` (value) where a pointer is required — read the compiler error.
3. Sketch `Publish(c *Course) error` that rejects nil / empty title.

## What the next lesson adds

**Lesson 12** — `defer`, `panic`, and `recover` used in the right places.
