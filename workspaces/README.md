
# 🚀 Go Workspaces Practice

Dự án này thực hành cơ chế **Go Workspaces** (giới thiệu từ bản Go 1.18). Cơ chế này cho phép làm việc với nhiều module cùng lúc trên môi trường local mà không cần sửa đổi file `go.mod` hay dùng lệnh `replace`.

## 📁 Cấu trúc thư mục
```text
workspaces/
├── go.work          <-- "Trái tim" của workspace
├── hello/           <-- Module chính (Main application)
│   ├── go.mod
│   └── main.go
└── utils/           <-- Module thư viện (Local library)
    ├── go.mod
    └── utils.go
```

## 🛠 Các bước đã thực hiện

### 1. Khởi tạo các Module
Tạo hai thư mục `hello` và `utils`, sau đó khởi tạo module cho từng cái:
```bash
go mod init example.com/hello
go mod init example.com/utils
```

### 2. Thiết lập Workspace
Tại thư mục gốc (`workspaces`), chạy lệnh để gộp các module vào một không gian làm việc chung:
```bash
go work init ./hello ./utils
```
Lệnh này tạo ra file `go.work`. Từ giờ, Go sẽ ưu tiên tìm các gói `import` trong các thư mục được liệt kê ở đây trước khi lên mạng tải.

### 3. Code & Chạy
Trong `hello/main.go`, ta có thể import trực tiếp:
```go
import "example.com/utils"
```
Mặc dù `example.com/utils` chưa hề được đẩy lên GitHub, Go vẫn nhận diện được nhờ file `go.work`.

## 💡 Tại sao dùng cái này?
* **Local Development:** Chỉnh sửa thư viện `utils` và thấy kết quả ngay lập tức bên `hello`.
* **Clean go.mod:** Không cần thêm các dòng `replace` tạm thời vào `go.mod`, giúp tránh lỗi khi commit code lên Git.
* **Tidy-friendly:** Lệnh `go mod tidy` sẽ không xóa các thư viện nội bộ này vì chúng đã được khai báo trong Workspace.

---

> **Tip:** Khi deploy lên server (Production), bạn thường sẽ xóa file `go.work` và đẩy các module lên Git theo đúng luồng quản lý phiên bản chuẩn.

---