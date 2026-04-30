## Demo: Chạy thẳng bằng Go vs Chạy qua Docker

### Bước 1 — Tạo project

Mở terminal (PowerShell hoặc CMD), chạy:

```powershell
mkdir demo-go
cd demo-go
go mod init demo
```

Tạo file `main.go`:
```powershell
notepad main.go
```

Paste vào:
```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, `{"message": "Hello từ Go!"}`)
    })

    fmt.Println("Server chạy tại http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

---

### Cách 1 — Chạy thẳng bằng Go (có source code)

```powershell
go run main.go
```

Mở browser: `http://localhost:8080/hello` → thấy response.

`Ctrl+C` để dừng.

---

### Cách 2 — Compile thành binary, xóa source, vẫn chạy được

```powershell
# Compile
go build -o mcp-server.exe main.go

# Xóa source đi
del main.go

# Kiểm tra thư mục — chỉ còn binary
dir

# Chạy binary (không cần Go, không cần source)
.\mcp-server.exe
```

Vào `http://localhost:8080/hello` → **vẫn chạy bình thường** dù đã xóa `main.go`.

`Ctrl+C` để dừng.

---

### Cách 3 — Build Docker image, chạy trong container

Tạo file `Dockerfile` (không có extension):
```powershell
notepad Dockerfile
```

Paste vào:
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /mcp-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /mcp-server .
EXPOSE 8080
CMD ["./mcp-server"]
```

```powershell
# Build image
docker build -t demo-go .

# Xem image size — nhỏ hơn Python rất nhiều
docker images demo-go

# Chạy container
docker run -p 8080:8080 demo-go
```

Vào `http://localhost:8080/hello` → chạy từ trong container.

---

## Điều bạn sẽ thấy rõ

| Bước | Quan sát |
|---|---|
| Cách 1 | Cần `go run` → cần Go toolchain |
| Cách 2 | Xóa source → binary vẫn chạy → **Go không cần thiết nữa** |
| Cách 3 | `docker images` → image chỉ ~15MB vs Python ~200MB |

Cách 2 chính là lý do Dockerfile Go dùng **multi-stage build** — stage 1 compile, stage 2 chỉ lấy binary, không mang theo Go toolchain vào production.