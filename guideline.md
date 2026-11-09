# Lesson 26 — Docker multi-stage Go binary

Branch: `lesson-26`  
Article: [Docker multi-stage: building a small Golang binary](https://tayjava.com/en/technologies/docker-multi-stage-building-a-small-golang-binary)

## Goal of this branch

Ship a two-stage `Dockerfile`: the builder compiles `./cmd/api` with `CGO_ENABLED=0`, `-trimpath`, and `-ldflags='-s -w'`; the runtime image is `gcr.io/distroless/static-debian12:nonroot` and only carries the binary at `/app`. Copy `go.mod` / `go.sum` before source so module downloads stay cached. Inject runtime config with `HTTP_ADDR` and `DATABASE_URL` (lesson 25).

## Changes vs `lesson-25`

| Add / change | Why |
|--------------|-----|
| `Dockerfile` | Multi-stage builder (`ARG GO_VERSION=1.22`) → distroless/static nonroot |
| `cmd/api` `-healthcheck` | Distroless has no curl; the binary probes `GET /healthz` on loopback |
| `.dockerignore` | Keep context small (no `.git`, `bin/`, tests, `.env`) |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` stays on the host; apply schema from the host or a migrate step (lesson 27):

```bash
psql "$DATABASE_URL" -f db-ddl.sql
```

Capacity column remains **`seats`**.

The image entrypoint is exec-form `["/app"]`, so the Go binary is PID 1 and receives `SIGTERM` directly. `HEALTHCHECK` and a Compose healthcheck both run `["/app", "-healthcheck"]` inside the container (no shell).

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
  -e DATABASE_URL='postgres://tayjava:tayjava@host.docker.internal:5432/courses?sslmode=disable' \
  tayjava/course-api:dev
```

On Linux you may need `--add-host=host.docker.internal:host-gateway`. After Postgres is reachable:

```bash
curl -s http://localhost:8080/healthz
curl -s http://localhost:8080/api/courses/1
docker inspect --format '{{.State.Health.Status}}' <container>
```

`-healthcheck` reads `HTTP_ADDR`, calls `http://127.0.0.1:<port>/healthz`, and exits 0 or 1. It runs before `config.Load()`, so it does not need `DATABASE_URL`. Compose can declare the same command:

```yaml
healthcheck:
  test: ["CMD", "/app", "-healthcheck"]
  interval: 10s
  timeout: 3s
  retries: 3
```

`go build ./...` and unit tests do not require Docker or a live database. Listen on `HTTP_ADDR` as `:8080` (all interfaces), not hard-coded `127.0.0.1`.

## Try yourself

1. Compare image size: single-stage `golang` vs this multi-stage tag.
2. Omit `DATABASE_URL` — the container should exit with the lesson-25 fail-fast message.
3. Confirm a curl-based healthcheck stays `unhealthy` on distroless; `["/app", "-healthcheck"]` does not.

## What the next lesson adds

**Lesson 27** — Intermediate lab: Course REST (`/api/courses`, `/api/courses/{id}`, `POST /api/enrollments`) + Docker Compose end-to-end.
