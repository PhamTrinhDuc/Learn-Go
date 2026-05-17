# Salon Chain Management - Backend API

Dự án Backend cho hệ thống Quản lý Chuỗi Salon/Spa (Salon Management System), được xây dựng bằng **Go (Golang)** theo kiến trúc **Clean Architecture**.

## 🚀 Công nghệ sử dụng (Tech Stack)

- **Ngôn ngữ:** Go (1.25+)
- **Web Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL (sử dụng driver `jackc/pgx/v5`)
- **Kiến trúc:** Clean Architecture (Domain - Repository - Usecase - Controller)
- **Xác thực (Authentication):** JWT (JSON Web Tokens) với RSA Keys (`golang-jwt/jwt/v5`)
- **Observability:** OpenTelemetry (Traces & Metrics) + OTLP Exporter
- **Khác:** UUID (`google/uuid`), Validator (`go-playground/validator`)

## 🏗️ Cấu trúc thư mục (Folder Structure)

Dự án được tổ chức theo chuẩn Clean Architecture để dễ dàng mở rộng và bảo trì:

- `api/`: Chứa định nghĩa Router và Middleware (ví dụ: `AuthMiddleware`, `TracingMiddleware`).
- `bootstrap/`: Khởi tạo các thành phần cốt lõi của ứng dụng như Database connection pool, OpenTelemetry, và Container (DI Container).
- `cmd/`: Điểm entrypoint của chương trình. Tệp `app.go` chứa hàm `main()`.
- `controller/`: Tầng Delivery/Handler, tiếp nhận request từ HTTP (Gin context), gọi Usecase và trả về response.
- `domain/`: Chứa các Model/Entity cốt lõi của hệ thống như `User`, `Branch`, `Stylist`, `Booking`, `Service`,...
- `usecase/`: Chứa business logic của ứng dụng (Tầng Use Case).
- `repository/`: Tầng giao tiếp với cơ sở dữ liệu PostgreSQL.
- `internal/`: Các package dùng chung hoặc tiện ích nội bộ như `auth`, `observability`, `utils`.
- `tmp/`: Chứa các file sinh ra trong quá trình chạy (ví dụ: RSA keys cho JWT).

## 🎯 1. Bài toán nghiệp vụ cốt lõi (Core Business Cases)

Hệ thống được thiết kế để số hóa toàn bộ quy trình vận hành của một chuỗi Salon từ lúc khách hàng tìm hiểu, đặt lịch, cho đến khi thực hiện dịch vụ và thanh toán.

### 1.1. Quản lý Đặt lịch thông minh (Smart Booking Engine)
Đây là "trái tim" của hệ thống (domain `Booking`). Khách hàng có thể đặt lịch qua nhiều nguồn: **Zalo, Web, AI Agent** hoặc **nhân viên đặt hộ (Manual)**.
- **Tìm khung giờ trống (Available Slots):** Hệ thống không chỉ lưu lịch hẹn, mà còn phải tính toán được *khoảng thời gian trống* dựa trên:
  1. Lịch làm việc trong tuần của Stylist (từ `StylistSchedule`).
  2. Thời lượng dự kiến của dịch vụ khách chọn (từ `Service.DurationMinutes`).
  3. Các booking hiện có của Stylist (để tránh double-booking).
- **Vòng đời của Booking:** `Pending` (Chờ xác nhận) ➔ `Confirmed` (Đã xác nhận) ➔ `Completed` (Đã hoàn thành).
- **Xử lý rủi ro:** 
  - *No-show Detection*: Tự động phát hiện các booking quá giờ mà khách không đến (GetNoShowBookings).
  - *Upcoming Reminders*: Truy xuất các lịch sắp tới để hệ thống Agent hoặc SMS gửi tin nhắn nhắc nhở khách.

### 1.2. Quản lý Thợ & Ca làm việc (Stylist & Scheduling)
Mỗi thợ làm tóc (`Stylist`) thuộc về một chi nhánh nhất định và có lịch làm việc (`StylistSchedule`) riêng.
- **Ca làm việc:** Cấu hình theo ngày trong tuần (Thứ 2 đến Chủ Nhật) với giờ bắt đầu (`StartTime`) và giờ kết thúc (`EndTime`).
- **Ràng buộc:** Khi khách hàng chọn một Stylist cụ thể, hệ thống sẽ cross-check lịch làm việc này trước khi cho phép Booking.

### 1.3. Cung cấp Dịch vụ (Service Catalog)
Danh mục dịch vụ (`Service`) được chia thành các nhóm (Categories): `Cut` (Cắt), `Color` (Nhuộm), `Treatment` (Phục hồi). 
- Thuộc tính quan trọng nhất là `EstimatedDuration` (Thời gian dự kiến). Khi khách đặt chuỗi nhiều dịch vụ (VD: Cắt + Nhuộm), hệ thống tự cộng dồn thời gian để block lịch của Stylist tương ứng.

### 1.4. Đơn hàng & Thanh toán (Orders & Billing)
Khi Booking chuyển sang trạng thái `Completed`, hệ thống sẽ khởi tạo Đơn hàng (`Orders`).
- **Thanh toán:** Hỗ trợ nhiều hình thức: Tiền mặt (`Cash`), Thẻ (`Card`), Trực tuyến (`Online`).
- Đơn hàng gắn liền với Khách hàng (`User`) và Chi nhánh (`Branch`) để phục vụ báo cáo doanh thu độc lập cho từng cơ sở.

### 1.5. Vận hành đa chi nhánh (Multi-Branch Management)
Tất cả các dữ liệu vận hành từ `Booking`, `Orders`, `Stylist`, đến `Inventory` (Kho hàng) đều được gắn chặt với một `BranchID`. Điều này cho phép:
- Quản lý tập trung trên một Database duy nhất (Multi-tenant logic ở level row).
- Báo cáo và phân quyền dữ liệu (Admin chi nhánh A chỉ xem được dữ liệu chi nhánh A).

---

*(Và các domain models khác đang phát triển: `Booking`, `Product`, `Inventory`, `Service`, `Orders`).*

## ⚙️ Biến môi trường (Environment Variables)

Hệ thống sử dụng các biến môi trường sau để cấu hình:

| Biến | Ý nghĩa | Mặc định |
|---|---|---|
| `PORT` | Cổng chạy HTTP Server | `8080` |
| `DB_HOST` | Địa chỉ máy chủ PostgreSQL | `localhost` |
| `DB_PORT` | Cổng kết nối PostgreSQL | `5433` |
| `DB_USER` | Tên đăng nhập DB | `mcp_user` |
| `DB_PASSWORD` | Mật khẩu DB | `mcp_password` |
| `DB_NAME` | Tên CSDL | `salon_chain` |
| `DB_SSLMODE` | Chế độ SSL cho DB | `disable` |
| `OTEL_SERVICE` | Tên service cho OpenTelemetry | `backend-server` |
| `OTEL_EXPORTER_OTLP_ENDPOINT`| Endpoint đẩy metrics/traces | `localhost:4318` |

## 🛠️ Hướng dẫn cài đặt và chạy thử

### 1. Yêu cầu hệ thống
- Đã cài đặt [Go 1.25+](https://go.dev/dl/).
- Đã cài đặt PostgreSQL (hoặc chạy qua Docker).
- (Tùy chọn) Redis, Jaeger / Prometheus nếu muốn dùng tracing/metrics.

### 2. Cài đặt các gói phụ thuộc
```bash
go mod tidy
```

### 3. Chạy ứng dụng
Vào thư mục chứa file `app.go` và chạy:
```bash
go run cmd/app.go
```
*Lưu ý: Ứng dụng sẽ tự động sinh cặp khóa RSA (`tmp/demo-keys/`) để test JWT và in một Test Token ra màn hình console khi khởi động thành công.*

### 4. API Endpoints
Base URL: `http://localhost:8080/api/v1`
Các nhóm API chính đã được thiết lập route:
- `/branch/*`
- `/user/*`
- `/stylist/*`
- `/stylist-schedule/*`
*(Tất cả API này hiện đang yêu cầu JWT Bearer Token, ngoại trừ các route public được cấu hình riêng).*
