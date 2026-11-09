# ADR-002 — Enrollment worker in Go beside a JVM Course API

Status: Accepted (lab)  
Date: lesson-42

## Context

TayJava may already run Course CRUD on Spring. Rewriting the monolith “because Go is fast” is rarely justified. Peak enrollment needs a bounded worker, clear pools, and a single binary ops story.

## Options

| Option | Pros | Cons |
|--------|------|------|
| A. Spring `@Async` | One language | Pool/tuning less explicit; fat deploy |
| B. Node BFF queue | JS team velocity | Weaker CPU-bound story; dual event loops |
| C. Go `cmd/worker` | Small binary, explicit pools, easy `-race` | Extra on-call language |

## Decision

Choose **C** for the **enrollment worker + metrics** bounded context only. Keep teacher-facing Course CRUD on the JVM if the team already owns it. Do not rewrite the whole domain.

## When not to choose Go

- Team has no Go on-call coverage.
- Problem is primarily heavy reporting / ORM-centric CRUD already healthy on Spring.
- Hiring pipeline cannot staff a second language.

## Consequences

- Shared runbook: `JWT_SECRET` / `DATABASE_URL` / metrics names must match the API lab.
- Polyglot is bounded: worker vs monolith, not a rewrite.
- Measure with `/metrics` and SLOs before declaring victory.
