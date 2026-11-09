# tayjava-go-lab

Lab code cho series **Golang từ Beginner tới Senior** (TayJava).

- Module: `github.com/tayjava/go-lab`
- Go: **1.22+**
- Mỗi bài học = một branch `lesson-NN` + file `guideline.md` giải thích thay đổi

Repo: https://github.com/luongquoctay87/tayjava-go-lab

## Bắt đầu nhanh

```bash
git clone https://github.com/luongquoctay87/tayjava-go-lab.git
cd tayjava-go-lab
git checkout lesson-02
go run .
```

Đọc `guideline.md` trên branch đang checkout trước khi sửa code.

## Quy ước branch

| Branch | Nội dung |
|--------|----------|
| `master` | Điểm xuất phát / README |
| `lesson-02` | Hello module — bài 2 |
| `lesson-03` … | Nâng dần theo series (sẽ thêm) |

Luôn học trên branch khớp số bài; đừng nhảy `master` rồi đoán thiếu file.

## Lệnh thường dùng

```bash
go version
go run .
go build -o bin/hello .
go test ./...
```
