# Lesson 9 — Idiomatic error handling (wrap, Is, As)

Branch: `lesson-09`  
Article: [Error handling idiomatic trong Golang](https://tayjava.com/vi/technologies/error-handling-idiomatic-trong-golang)

## Goal of this branch

Add sentinel `ErrNotFound`, wrap with `%w` in `Enroll`, and detect it with `errors.Is` from `cmd/errdemo`.

## Changes vs `lesson-08`

| Add / change | Why |
|--------------|-----|
| `course/store.go` | `MemoryStore`, `CourseStore`, `ErrNotFound`, `Enroll` |
| `cmd/errdemo/main.go` | Demo wrap + `errors.Is` |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-09
go run ./cmd/errdemo
go build ./...
```

Expected (shape):

```text
enroll error: enroll user 42 into course 9: course: not found
404-ish
Is with percent-v wrap: false
Is with percent-w wrap: true
```

## Try yourself

1. Temporarily change `Enroll` to wrap with `%v` and re-run — `errors.Is` fails.
2. Add a `ValidationError` type and try `errors.As`.
3. Save a course with ID 9 and confirm `Enroll` returns nil.

## What the next lesson adds

**Lesson 10** — table-driven tests with `go test` under `course/`.
