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
- Docker / Docker Compose optional for lessons 26–27
- Postgres from lesson 22 (`DATABASE_URL`)

## Quick start

```bash
git clone https://github.com/luongquoctay87/tayjava-go-lab.git
cd tayjava-go-lab
git checkout lesson-27   # or lesson-02 … lesson-26
cat guideline.md         # read before you edit
```

Lesson 2 hello:

```bash
git checkout lesson-02
go run .
# Hello, TayJava — Course lab
```

Lesson 27 Intermediate lab:

```bash
git checkout lesson-27
go test ./...
docker compose up --build
# curl -d '{"title":"Go Intermediate","seats":2}' -H 'Content-Type: application/json' http://localhost:8080/courses
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
| `lesson-18` | 18 | context timeout/cancel (`cmd/contextlab`) |
| `lesson-19` | 19 | net/http Course API (`cmd/api`) |
| `lesson-20` | 20 | JSON + logging/recovery middleware |
| `lesson-21` | 21 | Handler → service → memory repository |
| `lesson-22` | 22 | Postgres `database/sql` + `db-ddl.sql` |
| `lesson-23` | 23 | Validation + centralized HTTP errors |
| `lesson-24` | 24 | httptest table-driven HTTP tests |
| `lesson-25` | 25 | Env config / twelve-factor light |
| `lesson-26` | 26 | Docker multi-stage Course API image |
| `lesson-27` | 27 | Intermediate lab: Course REST + enrollment |
| `lesson-28` | 28 | Notify worker pool + fan-in / fan-out |
| `lesson-29` | 29 | errgroup pipeline + graceful HTTP shutdown |
| `lesson-30` | 30 | Race detector: unsafe vs safe course cache |
| `lesson-31` | 31 | Generics helpers (Unique / First) |
| `lesson-32` | 32 | Benchmarks + localhost pprof |
| `lesson-33` | 33 | slog JSON + request_id / course_id |
| `lesson-34` | 34 | Enroll counters + GET /metrics |
| `lesson-35` | 35 | Retry/backoff + Idempotency-Key |
| `lesson-36` | 36 | JWT Bearer middleware (`JWT_SECRET`) |
| `lesson-37` | 37 | gRPC vs REST checklist (light demo) |
| `lesson-38` | 38 | Enroll transaction + 3s deadline (`seats`) |

## Conventions

- **One lesson = one branch** `lesson-NN` (zero-padded: `02`, `03`, …).
- Every branch has **`guideline.md`** (read it first).
- Code on a branch must **run** with the commands in that guideline.
- **English only** in `.go` files, comments, commit-facing docs in this repo, and SQL comments.
- Postgres schema **`db-ddl.sql`**: present from **`lesson-22` onward**. Apply with `psql "$DATABASE_URL" -f db-ddl.sql`. Capacity column is **`seats`** (not `seats_left`).

## Layout (lesson 27)

```text
tayjava-go-lab/
├── README.md
├── guideline.md
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── .env.example
├── db-ddl.sql
├── go.mod
├── main.go
├── course/
├── internal/
│   ├── api/
│   ├── config/
│   ├── course/
│   └── seatledger/
├── cmd/ …
└── bin/
```

## Common commands

```bash
go version
go test ./...
go build -o bin/api ./cmd/api
docker compose up --build
git branch -a
git checkout lesson-27
# psql "$DATABASE_URL" -f db-ddl.sql
```

## Series & support

- Article roadmap: Golang series on TayJava (map starts at lesson 1; lab code from lesson 2).
- GitHub issues: lab won’t run / guideline drifts from the article.

## License

Lab content is for learning with the TayJava series. Keep attribution if you fork or reuse in an internal course.
