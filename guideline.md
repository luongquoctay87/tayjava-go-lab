# Lesson 26 — Docker multi-stage Go binary

Branch: `lesson-26`  
Article: [Docker multi-stage: build binary Golang nhỏ](https://tayjava.com/vi/technologies/docker-multi-stage-build-binary-golang-nho)

## Goal of this branch

Ship a multi-stage `Dockerfile`: builder compiles `./cmd/api` with `CGO_ENABLED=0`, runtime Distroless image only carries the binary. Layer-cache `go.mod` / `go.sum` before source. Inject runtime config via env (lesson 25).

## Changes vs `lesson-25`

| Add / change | Why |
|--------------|-----|
| `Dockerfile` | Multi-stage builder → distroless/static |
| `.dockerignore` | Keep context small (no `.git`, `bin/`, tests, `.env`) |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` stays in the image context if you need it on the host; apply schema from the host or a migrate step:

```bash
psql "$DATABASE_URL" -f db-ddl.sql
```

Capacity column remains **`seats`**.

## Layout

```text
Dockerfile
.dockerignore
.env.example
db-ddl.sql
cmd/api/main.go
internal/config/
```

## Run

```bash
git checkout lesson-26
go test ./...
go build -o bin/api ./cmd/api

docker build -t tayjava/course-api:dev .
docker run --rm -p 8080:8080 \
  -e HTTP_ADDR=:8080 \
  -e DATABASE_URL='postgres://tayjava:tayjava@host.docker.internal:5432/tayjava?sslmode=disable' \
  tayjava/course-api:dev
```

On Linux you may need `--add-host=host.docker.internal:host-gateway`. Hit `GET /healthz` after Postgres is reachable.

`go build ./...` and unit tests do not require Docker or a live database.

## Try yourself

1. Compare image size: single-stage `golang` vs this multi-stage tag.
2. Omit `DATABASE_URL` — container should exit with the lesson-25 fail-fast message.
3. Confirm `HTTP_ADDR` listens on `:8080` (all interfaces), not `127.0.0.1`.

## What the next lesson adds

**Lesson 27** — Intermediate lab: compose Course REST + enrollment + Docker Compose end-to-end.
