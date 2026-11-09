# Lesson 5 — Structs and methods (value vs pointer receivers)

Branch: `lesson-05`  
Article: [Struct và method — value vs pointer receiver](https://tayjava.com/vi/technologies/struct-va-method-value-vs-pointer-receiver)

## Goal of this branch

Define a `Course` struct, practice value vs pointer receivers, and see why mutation needs a pointer receiver.

## Changes vs `lesson-04`

| Add / change | Why |
|--------------|-----|
| `cmd/structdemo/main.go` | `Course` + `Publish` / `Remaining` / `Label` |
| `guideline.md` | This file |

## Run

```bash
git checkout lesson-05
go run ./cmd/structdemo
```

Expected (shape):

```text
after value Publish: false
after pointer Publish: {ID:1 Title:Go Seats:30 Published:true}
remaining(5): 25
label: 1:Go
{Title:Go Seats:1 Published:true}
Go Lab Go
```

## Try yourself

1. Flip `Publish` back to a value receiver and watch `Published` stay `false`.
2. Call `WithTitle` and confirm the original title is unchanged.
3. Add an embedded `Timestamps` struct and print a promoted field.

## What the next lesson adds

**Lesson 6** — small interfaces and composition: `cmd/ifacedemo` with `CourseStore` + `EnrollmentService`.
