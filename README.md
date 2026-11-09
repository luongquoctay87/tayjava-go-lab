# TayJava Go Lab

Companion source cho series **Golang từ Beginner tới Senior** của [TayJava](https://tayjava.com).

Một module Go xuyên suốt lộ trình — domain e-learning (`Course`, `Enrollment`, …). Mỗi bài học là một branch chạy được; file `guideline.md` trên branch đó giải thích code vừa thêm hoặc đổi gì.

| | |
|---|---|
| **Module** | `github.com/tayjava/go-lab` |
| **Go** | 1.22+ |
| **Remote** | https://github.com/luongquoctay87/tayjava-go-lab |

## Yêu cầu

- Go **1.22+** (`go version`)
- Git
- Terminal (macOS / Linux / Windows)

## Bắt đầu nhanh

```bash
git clone https://github.com/luongquoctay87/tayjava-go-lab.git
cd tayjava-go-lab
git checkout lesson-02
cat guideline.md    # đọc trước khi sửa
go run .
```

Kỳ vọng (bài 2):

```text
Hello, TayJava — Course lab
```

Build binary:

```bash
go build -o bin/hello .
./bin/hello
```

## Cách học với repo này

1. Mở bài viết trên TayJava (cùng số bài).
2. `git checkout lesson-NN` khớp số bài.
3. Đọc **`guideline.md`** — diff ý tưởng so với bài trước.
4. Chạy code, làm bài tập trong guideline / bài viết.
5. Sang bài sau: checkout branch kế tiếp (logic cộng dồn, không tạo project mới).

`master` chỉ là điểm neo / README. **Học trên `lesson-NN`**, không đoán thiếu file từ `master`.

## Branch hiện có

| Branch | Bài | Nội dung |
|--------|-----|----------|
| `master` | — | README, điểm xuất phát |
| `lesson-02` | 2 | Module + Hello (`go run` / `go build`) |

Các branch `lesson-03` … sẽ thêm theo series (biến/zero value → … → lab CLI → REST API → worker).

## Quy ước

- **Một bài = một branch** `lesson-NN` (số có zero-pad: `02`, `03`, …).
- Mỗi branch có **`guideline.md`** (bắt buộc đọc).
- Code trên branch phải **chạy được** với lệnh ghi trong guideline.
- Schema Postgres **`db-ddl.sql`**: xuất hiện từ **`lesson-22` trở đi** (Course / Enrollment); các branch trước đó không có file này.
- Identifier tiếng Anh; comment / guideline tiếng Việt.

## Cấu trúc (bài 2)

```text
tayjava-go-lab/
├── README.md
├── guideline.md      # đổi theo từng branch
├── go.mod
├── main.go           # entrypoint bài 2
├── .gitignore
└── bin/              # local build (gitignored)
```

Các bài sau sẽ thêm `cmd/`, `internal/`, … trên cùng module. `db-ddl.sql` chỉ có từ bài 22.

## Lệnh thường dùng

```bash
go version
go run .
go build -o bin/hello .
go test ./...
git branch -a
git checkout lesson-02
```

## Series & hỗ trợ

- Lộ trình bài viết: series Golang trên TayJava (bắt đầu từ bài 1 — bản đồ; lab code từ bài 2).
- Issue trên GitHub: lỗi chạy lab / lệch guideline so với bài viết.

## License

Nội dung lab phục vụ học tập kèm series TayJava. Giữ nguyên ghi chú nguồn khi fork hoặc dùng lại trong khóa nội bộ.
