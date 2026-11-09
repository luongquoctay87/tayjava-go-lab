# Lesson 23 — Validation and centralized HTTP errors

Branch: `lesson-23`  
Article: [Validation và xử lý lỗi HTTP tập trung trong Golang](https://tayjava.com/vi/technologies/validation-va-xu-ly-loi-http-tap-trung-trong-golang)

## Goal of this branch

Validate JSON at the HTTP edge (`DisallowUnknownFields` + simple rules) and map every domain error through one `ErrorBody` / `WriteError` / `MapDomainError` helper. Service still returns domain errors — never HTTP status codes.

## Changes vs `lesson-22`

| Add / change | Why |
|--------------|-----|
| `internal/api/errors.go` | `ErrorBody`, `WriteJSON`, `WriteError`, `MapDomainError` |
| `internal/api/course.go` | Decode/validate + shared error helpers |
| `internal/api/enrollment.go` | Decode/validate enroll body |
| `internal/api/handler_test.go` | Assert `VALIDATION` / `NO_SEATS` codes |
| `README.md` | Branch table through lesson 23 |
| `guideline.md` | This file |

`db-ddl.sql` is unchanged (`seats` column kept).

## Error contract

```json
{"code":"NO_SEATS","message":"no seats left"}
```

| Kind | Example | Status / code |
|------|---------|----------------|
| Input | `courseId: 0`, bad JSON | 400 `VALIDATION` |
| Business | no seats / already enrolled | 409 `NO_SEATS` / `CONFLICT` |
| Missing | unknown course | 404 `NOT_FOUND` |
| System | unexpected | 500 `INTERNAL` (log detail server-side) |

## Layout

```text
db-ddl.sql
internal/api/errors.go
internal/api/course.go
internal/api/enrollment.go
```

## Run

```bash
git checkout lesson-23
go test ./internal/api ./internal/course
go build -o bin/api ./cmd/api
./bin/api
```

```bash
curl -i -d '{"courseId":0,"userId":7}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
# 400 VALIDATION

curl -i -d '{"courseId":1,"userId":7}' -H 'Content-Type: application/json' \
  http://localhost:8080/api/enrollments
# enroll until seats=0 → 409 NO_SEATS
```

Apply schema when using Postgres:

```bash
psql "$DATABASE_URL" -f db-ddl.sql
```

## Try yourself

1. Send unknown JSON field (`course_id`) — expect `VALIDATION` / invalid json.
2. Force a non-domain error path and confirm the client only sees `INTERNAL`.
3. Keep service free of `net/http` imports.

## What the next lesson adds

**Lesson 24** — lock the error contract with table-driven `httptest` tests.
