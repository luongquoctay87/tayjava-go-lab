# Lesson 6 — Small interfaces and composition

Branch: `lesson-06`  
Article: [Interface nhỏ và composition trong Golang](https://tayjava.com/vi/technologies/interface-nho-va-composition-trong-golang)

## Goal of this branch

Define a small `CourseStore` interface, implement memory + stub stores, and inject them into `EnrollmentService` without changing the service.

## Changes vs `lesson-05`

| Add / change | Why |
|--------------|-----|
| `cmd/ifacedemo/main.go` | Interface + two store implementations |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-06
go run ./cmd/ifacedemo
```

Expected:

```text
Golang
Stub Golang
```

## Try yourself

1. Call `EnrollTitle` with a missing ID and print the error.
2. Add a third store (e.g. logging wrapper) without editing `EnrollmentService`.
3. Embed `CourseStore` into a larger interface that also has `Close()`.

## What the next lesson adds

**Lesson 7** — slices and maps: add `cmd/collections` with Course indexing by ID.
