# Lesson 4 — Functions, multiple return values, and scope

Branch: `lesson-04`  
Article: [Hàm, nhiều giá trị trả về và scope trong Golang](https://tayjava.com/vi/technologies/ham-nhieu-gia-tri-tra-ve-va-scope-trong-golang)

## Goal of this branch

Practice functions, `(T, error)`, pass-by-value, short-variable shadowing, and a first look at `defer`.

## Changes vs `lesson-03`

| Add / change | Why |
|--------------|-----|
| `cmd/funcdemo/main.go` | Demo helpers + exercise `canEnroll` |
| `guideline.md` | This file |

Previous demos (`cmd/zerodemo`, root `main.go`) remain runnable.

## Layout

```text
cmd/
  zerodemo/   # lesson 3
  funcdemo/   # lesson 4
```

## Run

```bash
git checkout lesson-04
go run ./cmd/funcdemo
```

Expected (shape):

```text
remaining: 20
canEnroll(30, 5) => ok=true err=<nil>
canEnroll(30, 40) => ok=false err=not enough seats
canEnroll(30, 0) => ok=false err=request must be positive
seats left after enroll: 25
after bump (pass by value): 10
enroll user 7 into course 101
done
```

## Try yourself

1. Change `canEnroll` error messages; re-run the three cases.
2. Add a named-return variant of `enroll` and compare with the explicit returns.
3. Observe that `bump` does not change the caller’s `n` (value semantics).

## What the next lesson adds

**Lesson 5** — structs and methods (value vs pointer receivers): add `cmd/structdemo` with a `Course` type.
