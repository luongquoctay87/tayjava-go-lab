# Lesson 34 — Basic metrics and tracing concepts

Branch: `lesson-34`  
Article: [Basic metrics and tracing in Golang — OpenTelemetry concepts](https://tayjava.com/vi/technologies/basic-metrics-and-tracing-in-golang-opentelemetry-concepts)

## Goal of this branch

Add enroll ok/fail counters and a temporary text `GET /metrics`. Sketch the idea of spans (logs already carry `request_id`); no Jaeger required.

## Changes vs `lesson-33`

| Add / change | Why |
|--------------|-----|
| `internal/metrics` | Atomic counters + `/metrics` text handler |
| `internal/api/enrollment.go` | Observe latency / outcome |
| `cmd/api/main.go` | Wire counters + `/metrics` |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged. Avoid high-cardinality labels like raw `course_id` on production metrics.

## Run

```bash
git checkout lesson-34
go test ./internal/metrics/...
go build -o bin/api ./cmd/api
# curl -s localhost:8080/metrics
```

## What the next lesson adds

**Lesson 35** — retry / backoff / idempotency for enroll side effects.
