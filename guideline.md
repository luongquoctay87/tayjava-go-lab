# Lesson 3 — Variables, types, and zero values

Branch: `lesson-03`  
Article: [Biến, kiểu dữ liệu và zero value trong Golang](https://tayjava.com/vi/technologies/bien-kieu-du-lieu-va-zero-value-trong-golang)

## Goal of this branch

Practice `var` vs `:=`, numeric conversion, and zero values on a Course-domain demo — same module as lesson 2.

## Changes vs `lesson-02`

| Add / change | Why |
|--------------|-----|
| `cmd/zerodemo/main.go` | Runnable demo for variables, casts, zero values |
| `guideline.md` | This file |

Root `main.go` from lesson 2 is unchanged (`go run .` still prints the Hello greeting).

Not yet: functions package, structs, interfaces, DB. Postgres schema (`db-ddl.sql`) starts at **lesson 22**.

## Layout

```text
tayjava-go-lab/
  go.mod
  main.go                 # lesson 2 entrypoint (kept)
  cmd/zerodemo/main.go    # lesson 3 demo
  guideline.md
  README.md
```

## Run

```bash
git checkout lesson-03
go run ./cmd/zerodemo
```

Expected (shape):

```text
course: Golang from zero seats: 25
price zero value: 0
remaining seats: 20
seats=0 title="" published=false course=<nil>
```

Also verify lesson 2 still works:

```bash
go run .
```

## Try yourself

1. Change `capacity` / `enrolled` and re-run; watch the cast.
2. Print zero values for `[]int` and `map[string]int` with `%#v`.
3. Try assigning without conversion (`capacity - enrolled`) and read the compiler error.

## What the next lesson adds

**Lesson 4** — functions, multiple return values, and scope: add `cmd/funcdemo` with `canEnroll` / `(T, error)` idioms.
