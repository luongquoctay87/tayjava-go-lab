# Lesson 31 — Generics where they pay off

Branch: `lesson-31`  
Article: [Generics in Golang — use them where they pay off](https://tayjava.com/vi/technologies/generics-in-golang-use-them-where-they-pay-off)

## Goal of this branch

Add small generic helpers (`Unique`, `First`, `MapIDs`). Do **not** genericize `EnrollmentService` — keep the Course domain concrete.

## Changes vs `lesson-30`

| Add / change | Why |
|--------------|-----|
| `internal/genericslab` | `Unique` / `First` / `MapIDs` + tests |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged. Note: we deliberately did not invent `Repository[T]`.

## Run

```bash
git checkout lesson-31
go test ./internal/genericslab/...
go build ./...
```

## What the next lesson adds

**Lesson 32** — benchmarks (`DescribeSlow`/`DescribeFast`) and `net/http/pprof` on `:6060`.
