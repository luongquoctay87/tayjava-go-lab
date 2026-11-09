# Lesson 36 — Basic JWT auth middleware

Branch: `lesson-36`  
Article: [Basic JWT auth middleware for a Course API in Golang](https://tayjava.com/vi/technologies/basic-jwt-auth-middleware-for-a-course-api-in-golang)

## Goal of this branch

Protect enroll routes with HS256 JWT Bearer middleware. Secret comes from `JWT_SECRET` (never commit a real secret). Lab login: `POST /login` with `{"userId":"7","password":"lab"}`.

## Changes vs `lesson-35`

| Add / change | Why |
|--------------|-----|
| `internal/auth` | Issue/Parse, Bearer, LoginHandler |
| `cmd/api/main.go` | Protect enroll; add `/login` |
| `.env.example` | `JWT_SECRET` |
| `go.mod` | `github.com/golang-jwt/jwt/v5` |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged (`idempotency_keys` still present from lesson 35).

## Run

```bash
git checkout lesson-36
export JWT_SECRET=lab-only-change-me
go test ./internal/auth/...
go build -o bin/api ./cmd/api
```

```bash
curl -s -d '{"userId":"7","password":"lab"}' -H 'Content-Type: application/json' localhost:8080/login
curl -s -d '{"userId":7}' -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  localhost:8080/courses/1/enrollments
```

## What the next lesson adds

**Lesson 37** — gRPC/protobuf concepts (light demo or skip note).  
**Lesson 38** — DB transactions + context deadlines on enroll (`seats` capacity).
