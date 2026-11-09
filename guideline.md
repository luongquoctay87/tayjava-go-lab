# Lesson 8 — Packages, modules, and export conventions

Branch: `lesson-08`  
Article: [Package, module (go.mod) và quy ước export](https://tayjava.com/vi/technologies/package-module-go-mod-va-quy-uoc-export)

## Goal of this branch

Extract a real `course` package with exported constructors/getters and call it from `cmd/coursedemo`.

## Changes vs `lesson-07`

| Add / change | Why |
|--------------|-----|
| `course/course.go` | Exported `Course`, `New`, `Level`; unexported `level` |
| `cmd/coursedemo/main.go` | Import the package and print title + level |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-08
go run ./cmd/coursedemo
go build ./...
```

Expected:

```text
Golang from zero beginner
```

Try accessing `c.level` from `main` — the compiler rejects it; use `c.Level()`.

## Try yourself

1. Add an unexported helper in `course/` and call it only from `New`.
2. Sketch a `course/store.go` file in the same package (no import cycle).
3. Confirm `go build ./...` still passes after edits.

## What the next lesson adds

**Lesson 9** — idiomatic errors: `ErrNotFound`, wrap with `%w`, `errors.Is` in `cmd/errdemo`.
