## Tại sao người ta chọn Compiled hay Interpreted?

### Compiled — ưu tiên hiệu năng

Khi thiết kế Go, C, Rust... mục tiêu là **chạy nhanh nhất có thể trên máy đích**. Để làm được vậy, compiler phải biết trước:

- Kiểu dữ liệu của mọi biến (int, string, struct...)
- Bộ nhớ cần cấp phát bao nhiêu
- CPU architecture là gì (x86, ARM...)

Vì biết tất cả trước, compiler tối ưu hóa được toàn bộ — sắp xếp lại lệnh, bỏ code thừa, tận dụng instruction set của CPU. Kết quả là binary chạy **gần với tốc độ tối đa của phần cứng**.

**Trade-off phải trả:**
- Phải khai báo kiểu dữ liệu tường minh → code dài hơn, verbose hơn
- Phải compile trước khi chạy → dev loop chậm hơn (sửa code → build → chạy)
- Binary chỉ chạy được trên OS/CPU đã compile cho → muốn chạy trên Windows và Linux phải build 2 lần

---

### Interpreted — ưu tiên linh hoạt

Khi thiết kế Python, Ruby... mục tiêu là **developer viết code nhanh nhất có thể**. Interpreter đọc code lúc chạy nên:

- Không cần khai báo kiểu — interpreter tự suy ra
- Sửa code → chạy ngay, không cần build
- Cùng file `.py` chạy được mọi nơi có Python

**Trade-off phải trả:**
- Chậm hơn compiled 10-100x vì phải dịch realtime
- Phải mang interpreter theo mọi nơi → Docker image to
- Lỗi kiểu dữ liệu chỉ phát hiện lúc chạy, không phải lúc build

---

### Không phải trắng đen hoàn toàn

Thực tế nhiều ngôn ngữ nằm ở giữa:

**Java, C#** — compile ra bytecode (không phải machine code), rồi JVM/CLR dịch tiếp lúc chạy. Lần đầu chạy chậm vì JVM "khởi động", sau đó JIT compiler tối ưu những đoạn code chạy nhiều → dần tiệm cận tốc độ compiled.

**Python** thực ra cũng compile ra `.pyc` (bytecode) trước, nhưng bytecode đó vẫn cần interpreter chạy — nên vẫn bị coi là interpreted.

**JavaScript** ban đầu thuần interpreted, nhưng V8 engine của Chrome dùng JIT — đo xem đoạn nào chạy nhiều rồi compile nóng đoạn đó → JS hiện tại nhanh hơn Python rất nhiều dù cùng "interpreted".

---

### Tóm lại — trade-off cốt lõi

| | Compiled | Interpreted |
|---|---|---|
| **Tốc độ chạy** | Nhanh | Chậm hơn |
| **Tốc độ viết code** | Chậm hơn | Nhanh |
| **Phát hiện lỗi** | Lúc build | Lúc chạy |
| **Portable** | Kém (phụ thuộc OS/CPU) | Tốt (chạy đâu có interpreter) |
| **Dùng cho** | Systems, infra, performance-critical | Scripting, data science, web backend |

Đó là lý do Go hợp cho infra tools (Redis, Prometheus, Docker đều viết bằng Go/C), còn Python hợp cho AI/data science — vì ở đó tốc độ viết code và thư viện quan trọng hơn hiệu năng thuần túy.