# Lesson 36 — Basic JWT auth middleware

Branch: `lesson-36`  
Article: [Basic JWT auth middleware for a Course API in Golang](https://tayjava.com/vi/technologies/basic-jwt-auth-middleware-for-a-course-api-in-golang)

## Goal of this branch

Protect enroll routes with HS256 JWT Bearer middleware. `JWT_SECRET` comes from the environment and must be at least 32 bytes; the process exits at startup when it is missing or shorter. The parser allows only HS256, requires `exp`, and allows 30 seconds of clock leeway. An empty or non-positive `sub` is rejected. The `Bearer` scheme is case-insensitive, and every 401 from auth sets `WWW-Authenticate`. No secret is hardcoded. Lab login: `POST /login` with `{"userId":7,"password":"lab"}`. Enroll uses the user id from the token, not the body.

## Changes vs `lesson-35`

| Add / change | Why |
|--------------|-----|
| `internal/auth` | Issue/Parse, Bearer, LoginHandler |
| `cmd/api/main.go` | Protect enroll; add `/login` |
| `.env.example` | `JWT_SECRET` |
| `go.mod` | `github.com/golang-jwt/jwt/v5` v5.2.1 (builds with Go 1.22) |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged (`idempotency_keys` still present from lesson 35).

## Run

```bash
git checkout lesson-36
export JWT_SECRET='dev-only-jwt-secret-change-me-32b!!'
go test ./internal/auth/...
go build -o bin/api ./cmd/api
```

```bash
curl -s -d '{"userId":7,"password":"lab"}' -H 'Content-Type: application/json' localhost:8080/login
curl -s -d '{"userId":7}' -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  localhost:8080/courses/1/enrollments
```

## What the next lesson adds

**Lesson 37** — gRPC/protobuf concepts (light demo or skip note).  
**Lesson 38** — DB transactions + context deadlines on enroll (`seats` capacity).
