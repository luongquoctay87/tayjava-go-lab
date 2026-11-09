# Lesson 37 — gRPC and protobuf (light demo)

Branch: `lesson-37`  
Article: [gRPC and protobuf in Golang — when to replace REST](https://tayjava.com/vi/technologies/grpc-and-protobuf-in-golang-when-to-replace-rest)

## Goal of this branch

Capture when gRPC pays off for internal services. This branch is a **light demo** (no `protoc` / generated stubs) so the lab stays runnable without the protobuf toolchain.

## Changes vs `lesson-36`

| Add / change | Why |
|--------------|-----|
| `cmd/grpcdemo` | Checklist demo: REST public / gRPC internal |
| `README.md` / `guideline.md` | This lesson |

`db-ddl.sql` unchanged. JWT enroll routes remain REST.

## Run

```bash
git checkout lesson-37
go run ./cmd/grpcdemo
go build ./...
```

## Why not full protobuf here

Generating `.pb.go` needs `protoc` + plugins on every learner machine. The article covers contracts; the Advanced lab (39) stays on HTTP + in-process workers.

## What the next lesson adds

**Lesson 38** — enroll inside a DB transaction with `context.WithTimeout`; keep **`seats`** as total capacity.
