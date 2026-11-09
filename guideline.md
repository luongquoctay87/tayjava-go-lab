# Lesson 12 — defer, panic, recover (use them in the right place)

Branch: `lesson-12`  
Article: [defer, panic, recover — dùng đúng chỗ](https://tayjava.com/vi/technologies/defer-panic-recover-dung-dung-cho)

## Goal of this branch

Use `defer` for cleanup (`Close`), demonstrate LIFO defer order, recover a panic once, then prefer returning `error` for expected failures.

## Changes vs `lesson-11`

| Add / change | Why |
|--------------|-----|
| `cmd/deferdemo/main.go` | `readCourses`, panic/recover, error refactor |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-12
go run ./cmd/deferdemo
```

Expected (shape):

```text
body
second defer
first defer
read bytes: ...
readCourses missing: open does-not-exist.json: ...
recovered enrollOrPanic: seats must not be negative
enroll error (no panic): seats must not be negative: -1
publish panicked: empty title
published: Go Beginner
```

## Try yourself

1. Remove `recover` around `enrollOrPanic(-1)` and watch the process crash.
2. Confirm missing files return an error from `readCourses` — no panic.
3. Prefer `error` returns for business rules; reserve panic for true programmer bugs.

## What the next lesson adds

**Lesson 13** — beginner lab CLI: `cmd/course-cli` + `internal/course` + JSON file store.
