# Go Context - Giải thích cho người mới bắt đầu

## Context là gì?

Hãy tưởng tượng bạn đặt đồ ăn online. Khi bạn hủy đơn hàng, hệ thống cần **thông báo cho tất cả các bên liên quan** (nhà bếp, shipper, kho...) dừng xử lý đơn đó.

**Context trong Go hoạt động y như vậy** — nó là một cơ chế để truyền tín hiệu hủy, deadline, và dữ liệu xuyên suốt các goroutine liên quan đến một request.

---

## Interface Context có 4 phương thức

```go
type Context interface {
    Done()     <-chan struct{}       // Kênh báo hiệu "hãy dừng lại"
    Deadline() (time.Time, bool)    // Thời điểm context sẽ bị hủy
    Err()      error                // Lý do context bị hủy
    Value(key any) any              // Lấy giá trị được gắn vào context
}
```

| Phương thức | Ý nghĩa đơn giản |
|---|---|
| `Done()` | Trả về channel — khi channel này đóng lại, có nghĩa là "hãy dừng công việc" |
| `Deadline()` | Hỏi "context này sẽ hết hạn lúc mấy giờ?" |
| `Err()` | Hỏi "tại sao context bị hủy?" |
| `Value()` | Lấy dữ liệu được đính kèm trong context |

---

## Tạo Context như thế nào?

### 1. `context.Background()` — Điểm xuất phát
```go
ctx := context.Background()
```
Đây là context **gốc rễ**, không bao giờ bị hủy. Dùng trong `main()`, tests, hoặc làm nền tảng để tạo context khác.

### 2. `context.TODO()` — Tạm thời chưa biết dùng gì
```go
ctx := context.TODO()
```
Dùng khi **chưa chắc** nên dùng context nào — đây là dấu hiệu cho team biết "cần bổ sung sau".

---

## 3 cách dùng Context phổ biến

### ① Truyền dữ liệu — `WithValue`

```go
func main() {
    ctx := context.Background()

    // Gắn userID vào context
    ctx = context.WithValue(ctx, "userID", "u-123")

    handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
    userID := ctx.Value("userID")
    fmt.Println("Đang xử lý cho user:", userID) // u-123
}
```

> ⚠️ **Lưu ý:** Chỉ dùng `WithValue` cho dữ liệu phụ trợ (request ID, user info...), **không dùng** cho tham số quan trọng của hàm.

---

### ② Hủy thủ công — `WithCancel`

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go doWork(ctx)

    time.Sleep(2 * time.Second)
    cancel() // Ra lệnh dừng goroutine
}

func doWork(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Dừng lại:", ctx.Err()) // context canceled
            return
        default:
            fmt.Println("Đang làm việc...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

> 💡 Luôn gọi `cancel()` khi xong việc để giải phóng tài nguyên (thường dùng `defer cancel()`).

---

### ③ Tự động hủy theo thời gian — `WithTimeout` / `WithDeadline`

```go
func main() {
    // Tự động hủy sau 3 giây
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    select {
    case <-time.After(5 * time.Second):
        fmt.Println("Hoàn thành!")  // Không bao giờ chạy tới đây

    case <-ctx.Done():
        fmt.Println("Hết giờ:", ctx.Err()) // context deadline exceeded
    }
}
```

**Sự khác biệt giữa hai hàm:**

| Hàm | Cách dùng |
|---|---|
| `WithTimeout(ctx, 3*time.Second)` | Hủy sau **N giây** kể từ bây giờ |
| `WithDeadline(ctx, time.Now().Add(...))` | Hủy vào **một thời điểm cụ thể** |

Thực ra `WithTimeout` chỉ là wrapper của `WithDeadline`:
```go
// Bên trong WithTimeout thực ra là:
return WithDeadline(parent, time.Now().Add(timeout))
```

---

## Ví dụ thực tế: HTTP Server

```go
func handleRequest(w http.ResponseWriter, req *http.Request) {
    ctx := req.Context() // Mỗi HTTP request đã có sẵn context

    select {
    case <-time.After(5 * time.Second):
        fmt.Fprintf(w, "Kết quả từ server")

    case <-ctx.Done(): // Client đóng kết nối trước khi server xong
        fmt.Println("Client đã hủy request:", ctx.Err())
    }
}
```

Khi client bấm **Ctrl+C** hủy request, server **biết ngay** và dừng xử lý — tránh lãng phí tài nguyên.

---

## Tóm tắt "khi nào dùng cái gì"

```
Cần điểm khởi đầu          → context.Background()
Chưa biết dùng gì           → context.TODO()
Truyền dữ liệu theo request → WithValue()
Cho phép hủy thủ công       → WithCancel()
Giới hạn thời gian chạy     → WithTimeout() hoặc WithDeadline()
```

---

## Quy tắc vàng khi dùng Context

> 1. **Luôn truyền context là tham số đầu tiên** của hàm, đặt tên là `ctx`
> 2. **Luôn gọi `cancel()`** — dùng `defer cancel()` ngay sau khi tạo
> 3. **Không lưu context vào struct** — truyền trực tiếp qua hàm
> 4. **Không truyền `nil` context** — dùng `context.TODO()` nếu chưa chắc