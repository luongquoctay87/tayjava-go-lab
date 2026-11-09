# Lesson 2 — What is Go? Install the toolchain and run your first program

Branch: `lesson-02`  
Article: [Go là gì? Cài toolchain và chạy chương trình Golang đầu tiên](https://tayjava.com/vi/technologies/go-la-gi-cai-toolchain-va-chay-chuong-trinh-golang-dau-tien)

## Goal of this branch

First runnable **lab module**: print a greeting; understand `go run` vs `go build`.

## Changes vs `master` / previous

| Add / change | Why |
|--------------|-----|
| `go.mod` | Module `github.com/tayjava/go-lab` (module mode, Go 1.22+) |
| `main.go` | `package main` + `Hello, TayJava — Course lab` |
| `.gitignore` | Ignore `bin/` and local artifacts |
| `README.md` | How to clone and check out by lesson |
| `guideline.md` | This file — read before editing code |

Not yet: `cmd/`, `internal/`, HTTP, DB. Postgres schema (`db-ddl.sql`) starts at **lesson 22**.

## Layout

```text
tayjava-go-lab/
  go.mod
  main.go          ← lesson 2 entrypoint (repo root)
  guideline.md
  README.md
  .gitignore
  bin/             ← created by go build (gitignored)
```

## Run

Requires Go **1.22+** (`go version`).

```bash
git checkout lesson-02
go run .
```

Expected:

```text
Hello, TayJava — Course lab
```

Build a binary:

```bash
go build -o bin/hello .
./bin/hello
```

## Try yourself

1. Read `main.go`, change the printed string, run `go run .` again.
2. Compare when files appear under `bin/` after `go build`.
3. Commit experiments on your own branch (push optional).

## What the next lesson adds

**Lesson 3** — variables, types, zero values: add a demo (e.g. `cmd/zerodemo`) on the same module; do not create a new repo.
