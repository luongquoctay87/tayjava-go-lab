# Lesson 10 — Unit tests and table-driven tests

Branch: `lesson-10`  
Article: [Unit test và table-driven tests với go test](https://tayjava.com/vi/technologies/unit-test-va-table-driven-tests-voi-go-test)

## Goal of this branch

Export `RemainingSeats`, write black-box table-driven tests, and assert `ErrNotFound` with `errors.Is`.

## Changes vs `lesson-09`

| Add / change | Why |
|--------------|-----|
| `course/seats.go` | `RemainingSeats` + `ErrInvalidEnrollment` |
| `course/seats_test.go` | `package course_test` table-driven tests |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-10
go test ./course
go test -v ./course
go test -run RemainingSeats ./course
go test ./...
```

Expected: all tests PASS. Verbose output shows subtests such as `TestRemainingSeats/over`.

## Try yourself

1. Change one `want` value and re-run — the failure message includes the case name.
2. Add a case for `enrolled < 0`.
3. Prefer `errors.Is` for sentinel errors in assertions.

## What the next lesson adds

**Lesson 11** — pointers: `reserve` / `reserveCopy` on `Course` seats.
