# Lesson 13 — Beginner lab: Course CLI + JSON file

Branch: `lesson-13`  
Article: [Lab Beginner Golang: CLI quản lý Course (JSON file)](https://tayjava.com/vi/technologies/lab-beginner-golang-cli-quan-ly-course-json-file)

## Goal of this branch

Ship a small Course management CLI that persists to `data/courses.json`, using `internal/course` so other packages outside the module tree cannot import the domain package by accident.

## Changes vs `lesson-12`

| Add / change | Why |
|--------------|-----|
| `internal/course/model.go` | JSON `Course` with exported `Level` |
| `internal/course/store.go` | `Load` / `Save` / `Add` / `Get` |
| `internal/course/store_test.go` | Duplicate `Add` + missing-file tests |
| `cmd/course-cli/main.go` | `add` / `list` / `get` subcommands |
| `.gitignore` | Ignore generated `data/*.json` |
| `data/.gitkeep` | Keep the data directory in git |
| `README.md` | Branch table through lesson 13 |
| `guideline.md` | This file |

Earlier packages stay: root `main.go`, `cmd/*` demos, and top-level `course/` (lessons 8–11). The CLI uses **`internal/course`** as the article’s capstone layout.

Not yet: HTTP API, Postgres, `db-ddl.sql` (lesson 22+).

## Layout

```text
tayjava-go-lab/
  cmd/course-cli/main.go
  internal/course/
    model.go
    store.go
    store_test.go
  data/                 # courses.json created on first add (gitignored)
  course/               # earlier package demos (kept)
  cmd/…                 # lessons 3–12 demos
```

## Run

```bash
git checkout lesson-13
mkdir -p data
go test ./internal/course ./course
go run ./cmd/course-cli add -id 1 -title "Golang from zero" -seats 30
go run ./cmd/course-cli add -id 2 -title "Go Concurrency" -seats 20 -level intermediate
go run ./cmd/course-cli list
go run ./cmd/course-cli get -id 1
```

Expected (shape):

```text
added course 1
added course 2
1	Golang from zero	seats=30	beginner
2	Go Concurrency	seats=20	intermediate
&{ID:1 Title:Golang from zero Seats:30 Level:beginner}
```

Duplicate id exits non-zero with a message on stderr:

```bash
go run ./cmd/course-cli add -id 1 -title "dup" -seats 1
# course id 1 already exists
```

## Try yourself

1. Delete `data/courses.json` and run `list` — empty output, no crash.
2. Add an `update -id N -seats X` subcommand.
3. Reject empty titles in `Add`.

## What the next lesson adds

**Lesson 14+** — concurrency / intermediate track (see series roadmap). Schema file `db-ddl.sql` still waits until **lesson 22**.
