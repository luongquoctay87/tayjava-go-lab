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
git checkout lesson-02
cat guideline.md    # read before you edit
go run .
```

Expected (lesson 2):

```text
Hello, TayJava — Course lab
```

Build a binary:

```bash
go build -o bin/hello .
./bin/hello
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

Branches `lesson-03` … will follow the series (zero values → … → CLI lab → REST API → worker).

## Conventions

- **One lesson = one branch** `lesson-NN` (zero-padded: `02`, `03`, …).
- Every branch has **`guideline.md`** (read it first).
- Code on a branch must **run** with the commands in that guideline.
- **English only** in `.go` files, comments, commit-facing docs in this repo, and SQL comments.
- Postgres schema **`db-ddl.sql`**: introduced from **`lesson-22` onward** (Course / Enrollment); earlier branches do not include it.

## Layout (lesson 2)

```text
tayjava-go-lab/
├── README.md
├── guideline.md      # per-branch
├── go.mod
├── main.go           # lesson 2 entrypoint
├── .gitignore
└── bin/              # local build (gitignored)
```

Later lessons add `cmd/`, `internal/`, … on the same module. `db-ddl.sql` appears from lesson 22.

## Common commands

```bash
go version
go run .
go build -o bin/hello .
go test ./...
git branch -a
git checkout lesson-02
```

## Series & support

- Article roadmap: Golang series on TayJava (map starts at lesson 1; lab code from lesson 2).
- GitHub issues: lab won’t run / guideline drifts from the article.

## License

Lab content is for learning with the TayJava series. Keep attribution if you fork or reuse in an internal course.
