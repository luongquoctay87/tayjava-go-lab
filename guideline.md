# Lesson 42 — Go vs JVM vs Node (worker binary)

Branch: `lesson-42` (from `lesson-40`; **skips lesson-41**)  
Article: [Go vs JVM vs Node — picking Golang trade-offs](https://tayjava.com/vi/technologies/go-vs-jvm-vs-node-picking-golang-trade-offs-for-your-team)

## Why lesson-41 is skipped

Lesson 41 is a **system-design** brief (sync/async/hybrid ADR, peak enrollment pitch). No new production code is required beyond what lesson 39 already demonstrates (bounded queue + 503). Read the article for the half-page ADR; this lab path continues at 42 with a concrete `cmd/worker`.

## Goal of this branch

Ship a small standalone `cmd/worker` and an ADR scoring Go vs JVM vs Node for the enrollment worker — not a Course CRUD rewrite.

## Changes vs `lesson-40`

| Add / change | Why |
|--------------|-----|
| `cmd/worker` | Standalone Go worker binary |
| `docs/adr-002-go-vs-jvm-node.md` | Trade-off ADR |
| `README.md` / `guideline.md` | Skip-41 note + this lesson |

`db-ddl.sql` unchanged.

## Run

```bash
git checkout lesson-42
export DATABASE_URL=postgres://tayjava:tayjava@localhost:5432/tayjava?sslmode=disable
go build -o bin/worker ./cmd/worker
go build ./...
cat docs/adr-002-go-vs-jvm-node.md
```

`go run ./cmd/worker` needs Postgres + a course id=1 (Compose from lesson 27). SIGTERM drains the queue.
