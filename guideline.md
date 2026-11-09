# Lesson 7 — Slices and maps

Branch: `lesson-07`  
Article: [Slice và map trong Golang](https://tayjava.com/vi/technologies/slice-va-map-trong-golang)

## Goal of this branch

Build Course lists with slices, index them with maps, and see why a nil map panics on write.

## Changes vs `lesson-06`

| Add / change | Why |
|--------------|-----|
| `cmd/collections/main.go` | Slice/map demos + Course index exercise |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-07
go run ./cmd/collections
```

Expected (shape):

```text
3 Golang
[101 102] 2 8
nil slice: true empty slice: false
len: 3 cap: ...
lookup 2: {ID:2 Title:Java Seats:15} ok=true
nil map write panicked (expected): assignment to entry in nil map
fixed map len: 1
enrolled course: 101
not enrolled
```

(`cap` after `append` may be 3 or larger depending on growth.)

## Try yourself

1. Remove the `recover` in `demoNilMapPanic` and watch the process crash.
2. Always initialize maps with `make` or a literal before writing.
3. Compare `var s []Course` vs `s := []Course{}` for JSON / `nil` checks later.

## What the next lesson adds

**Lesson 8** — packages, modules, and export: extract a real `course` package + `cmd/coursedemo`.
