# ADR-001 — Keep enrollment inside this module (not a separate Go module)

Status: Accepted (lab)  
Date: lesson-40

## Context

Lesson 39 grew HTTP handlers, an in-process worker, metrics, and Postgres transactions. Splitting `enrollment` into its own module (`github.com/tayjava/enrollment`) would force versioned releases for every JSON/SQL tweak.

## Decision

Keep Course + Enrollment under `github.com/tayjava/go-lab` with implementations in `internal/`. Public surface stays small (`course` package helpers for early lessons; API lives in `cmd/` + `internal/`).

## Consequences

- Refactors of SQL/DTO stay cheap (no `/v2` path yet).
- Consumers outside the module cannot import `internal/...` (enforced by the Go toolchain).
- If a second binary needs shared types later, extract a tiny `pkg/` or a v1 module with a stable DTO — not the Postgres repo.

## Compatibility note

Adding `Course.subtitle` is additive JSON. Removing or renaming `seats` would be a breaking change and needs a process (changelog / v2).
