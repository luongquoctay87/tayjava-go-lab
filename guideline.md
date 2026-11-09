# Bài 2 — Go là gì? Cài toolchain và chạy chương trình đầu tiên

Branch: `lesson-02`  
Bài viết: [Go là gì? Cài toolchain và chạy chương trình Golang đầu tiên](https://tayjava.com/vi/technologies/go-la-gi-cai-toolchain-va-chay-chuong-trinh-golang-dau-tien)

## Mục tiêu branch này

Lần đầu có **module lab** chạy được: in lời chào, hiểu `go run` vs `go build`.

## Thay đổi so với `master` / trước đó

| Thêm / đổi | Giải thích |
|------------|------------|
| `go.mod` | Module `github.com/tayjava/go-lab` (module mode, Go 1.22+) |
| `main.go` | `package main` + `Hello, TayJava — Course lab` |
| `.gitignore` | Bỏ qua `bin/` và artifact local |
| `README.md` | Cách clone và checkout theo bài |
| `guideline.md` | File này — đọc trước khi sửa code |

Chưa có: `cmd/`, `internal/`, HTTP, DB. Schema Postgres (`db-ddl.sql`) sẽ có từ **bài 22**.

## Cấu trúc hiện tại

```text
tayjava-go-lab/
  go.mod
  main.go          ← entrypoint bài 2 (root)
  guideline.md
  README.md
  .gitignore
  bin/             ← tạo khi go build (đã gitignore)
```

## Chạy thử

Yêu cầu: Go **1.22+** (`go version`).

```bash
git checkout lesson-02
go run .
```

Kỳ vọng:

```text
Hello, TayJava — Course lab
```

Build binary:

```bash
go build -o bin/hello .
./bin/hello
```

## Việc bạn nên tự làm (học)

1. Đọc `main.go`, đổi chuỗi in, chạy lại `go run .`.
2. So sánh thời điểm tạo file trong `bin/` khi `go build`.
3. Commit thay đổi thử trên branch riêng của bạn (không bắt buộc push).

## Bài sau sẽ thêm gì?

**Bài 3** — biến, kiểu, zero value: thêm demo (ví dụ `cmd/zerodemo`) trên cùng module, không tạo repo mới.
