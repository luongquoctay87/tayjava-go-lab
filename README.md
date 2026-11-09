# TayJava Go Lab

Companion source for the **Golang from Beginner to Senior** series by [TayJava](https://tayjava.com).

One Go module across the whole path — e-learning domain (`Course`, `Enrollment`, …). Each lesson is a runnable git branch; that branch’s `guideline.md` explains what was added or changed.

| | |
|---|---|
| **Module** | `github.com/tayjava/go-lab` |
| **Go** | 1.22+ |
| **Remote** | https://github.com/luongquoctay87/tayjava-go-lab |
| **Language** | **English** for all source, comments, SQL comments, `README.md`, and `guideline.md` |

## Requirements

- Go **1.22+** (`go version`)
- Git
- Terminal (macOS / Linux / Windows)

## Quick start

```bash
git clone https://github.com/luongquoctay87/tayjava-go-lab.git
cd tayjava-go-lab
git checkout lesson-17   # or lesson-02 … lesson-16
cat guideline.md         # read before you edit
```

Lesson 2 hello:

```bash
git checkout lesson-02
go run .
# Hello, TayJava — Course lab
```

Lesson 13 Course CLI:

```bash
git checkout lesson-13
go run ./cmd/course-cli add -id 1 -title "Golang from zero" -seats 30
go run ./cmd/course-cli list
```

Concurrency demos (14–17):

```bash
git checkout lesson-17
go run ./cmd/select-lab
go test -race ./internal/seatledger
go run ./cmd/synclab
```

## How to learn with this repo

1. Open the matching TayJava article (same lesson number).
2. `git checkout lesson-NN` for that lesson.
3. Read **`guideline.md`** — conceptual diff vs the previous lesson.
4. Run the code; do the exercises in the guideline / article.
5. Move on: check out the next branch (cumulative lab — no new project).

`master` is only an anchor / README. **Learn on `lesson-NN`**; do not guess missing files from `master`.

## Branches available

| Branch | Lesson | Contents |
|--------|--------|----------|
| `master` | — | README, starting point |
| `lesson-02` | 2 | Module + Hello (`go run` / `go build`) |
| `lesson-03` | 3 | Variables / zero values (`cmd/zerodemo`) |
| `lesson-04` | 4 | Functions / multi-return (`cmd/funcdemo`) |
| `lesson-05` | 5 | Struct + methods (`cmd/structdemo`) |
| `lesson-06` | 6 | Small interfaces (`cmd/ifacedemo`) |
| `lesson-07` | 7 | Slice + map (`cmd/collections`) |
| `lesson-08` | 8 | Package export (`course`, `cmd/coursedemo`) |
| `lesson-09` | 9 | Errors wrap / Is (`cmd/errdemo`) |
| `lesson-10` | 10 | Table-driven tests (`course/*_test.go`) |
| `lesson-11` | 11 | Pointers (`cmd/ptrdemo`) |
| `lesson-12` | 12 | defer / panic / recover (`cmd/deferdemo`) |
| `lesson-13` | 13 | Course CLI + JSON (`cmd/course-cli`, `internal/course`) |
| `lesson-14` | 14 | Goroutines (`cmd/goroutine-lab`) |
| `lesson-15` | 15 | Channels (`cmd/channel-lab`) |
| `lesson-16` | 16 | select / buffered / close (`cmd/select-lab`) |
| `lesson-17` | 17 | WaitGroup + Mutex (`internal/seatledger`, `cmd/synclab`) |

Later branches continue the series (context → REST API → worker). Postgres schema (`db-ddl.sql`) from **lesson 22**.

## Conventions

- **One lesson = one branch** `lesson-NN` (zero-padded: `02`, `03`, …).
- Every branch has **`guideline.md`** (read it first).
- Code on a branch must **run** with the commands in that guideline.
- **English only** in `.go` files, comments, commit-facing docs in this repo, and SQL comments.
- Postgres schema **`db-ddl.sql`**: introduced from **`lesson-22` onward** (Course / Enrollment); earlier branches do not include it.

## Layout (lesson 17)

```text
tayjava-go-lab/
├── README.md
├── guideline.md
├── go.mod
├── main.go
├── course/
├── internal/
│   ├── course/            # CLI JSON store
│   └── seatledger/        # lesson 17 Mutex ledger
├── cmd/
│   ├── course-cli/
│   ├── goroutine-lab/ … select-lab/
│   └── synclab/           # lesson 17 demo
├── data/
└── bin/
```

`db-ddl.sql` appears from lesson 22.

## Common commands

```bash
go version
go run .
go run ./cmd/course-cli list
go run ./cmd/goroutine-lab
go run ./cmd/channel-lab
go run ./cmd/select-lab
go run ./cmd/synclab
go test -race ./internal/seatledger
go test ./...
git branch -a
git checkout lesson-17
```

## Series & support

- Article roadmap: Golang series on TayJava (map starts at lesson 1; lab code from lesson 2).
- GitHub issues: lab won’t run / guideline drifts from the article.

## License

Lab content is for learning with the TayJava series. Keep attribution if you fork or reuse in an internal course.
