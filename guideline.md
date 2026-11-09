# Lesson 40 — Package design and backward compatibility

Branch: `lesson-40`  
Article: [Golang package design and backward compatibility](https://tayjava.com/vi/technologies/golang-package-design-and-backward-compatibility)

## Goal of this branch

Tighten package boundaries: keep SQL/HTTP under `internal`, document exported symbols, add an additive JSON field (`subtitle`), and record why enrollment is not a separate module yet.

## Changes vs `lesson-39`

| Add / change | Why |
|--------------|-----|
| `docs/adr-001-enrollment-module.md` | Short ADR: stay mono-module |
| `internal/course/model.go` | Additive `subtitle` JSON field |
| `README.md` / `guideline.md` | Export inventory + this lesson |

`db-ddl.sql` unchanged (subtitle is API-only for now; persistence can ignore omitempty).

## Export inventory (who imports?)

| Symbol / package | Importers |
|------------------|-----------|
| `course` (root) | Early demos (`cmd/coursedemo`, benches) |
| `internal/course` | `cmd/api`, `internal/api`, `internal/enrollment` |
| `internal/enrollment` | `cmd/api` only |
| `internal/auth` | `cmd/api` |
| `internal/notify` | `cmd/api`, `internal/api` |

Rule of thumb: prefer `internal/` for anything with SQL or HTTP wiring.

## Run

```bash
git checkout lesson-40
go test ./...
go build ./...
cat docs/adr-001-enrollment-module.md
```

## Skipped next

**Lesson 41** (system design ADR / peak enrollment) is documentation-heavy; skipped in this lab path — see article for the sync/async brief. Continue to **lesson 42** (`cmd/worker` + Go vs JVM vs Node ADR).
