# Lesson 32 — Benchmarks and pprof

Branch: `lesson-32`  
Article: [Benchmarks and pprof in Golang](https://tayjava.com/vi/technologies/benchmarks-and-pprof-in-golang-find-cpu-and-heap-hotspots)

## Goal of this branch

Compare `DescribeSlow` vs `DescribeFast` with benchmarks (`-benchmem`) and attach `net/http/pprof` on `127.0.0.1:6060`.

## Changes vs `lesson-31`

| Add / change | Why |
|--------------|-----|
| `course/describe.go` | Slow concat vs `strings.Builder` |
| `course/describe_test.go` | Benchmarks + equality check |
| `cmd/api/main.go` | Localhost pprof server |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged.

## Run

```bash
git checkout lesson-32
go test ./course -bench=Describe -benchmem
go build -o bin/api ./cmd/api
```

Profile (API running):

```bash
go tool pprof -http=:8081 http://127.0.0.1:6060/debug/pprof/profile?seconds=15
```

## What the next lesson adds

**Lesson 33** — structured logging with `log/slog` (JSON, `course_id` / `request_id`).
