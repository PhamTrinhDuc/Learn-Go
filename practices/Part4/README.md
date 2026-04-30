## 1. Bài tập 1: The Pipeline Builder (Functional Programming)
**Yêu cầu:** Thiết kế một hàm `Pipeline` nhận vào một danh sách các hàm (slice of functions). Mỗi hàm trong slice có dạng `func(int) int`. Hàm `Pipeline` phải trả về một hàm mới. Khi gọi hàm mới này với một số nguyên $x$, nó sẽ thực thi lần lượt các hàm trong pipeline theo thứ tự từ trái sang phải, lấy kết quả của hàm trước làm đầu vào cho hàm sau.

* **Đặc điểm:** Sử dụng Higher-order functions.
* **Input ví dụ:** `Pipeline(addOne, square, double)(2)` 
* **Giải thích:** $((2 + 1)^2) \times 2 = 18$.

## 2. Bài tập 2: Memoization Decorator (Closures)
**Yêu cầu:** Viết một hàm `Memoize` nhận vào một hàm tính toán "đắt đỏ" `fn func(int) int`. Hàm `Memoize` trả về một hàm có cùng signature nhưng có khả năng lưu trữ kết quả (cache). Nếu gọi lại với cùng một đối số, nó phải trả về kết quả từ bộ nhớ thay vì tính toán lại.

* **Đặc điểm:** Áp dụng Closures và Map để quản lý state bên trong hàm.
* **Thách thức:** Đảm bảo tính thread-safe nếu bạn muốn nâng cấp lên mức độ Expert (sử dụng `sync.Mutex`).

## 🟢 Bài 1: Middleware Chain (Thực tế: HTTP Frameworks)
**Bối cảnh:** Bạn đang xây dựng một Web Framework siêu nhẹ. Bạn cần một cơ chế để các "Middleware" (như Logging, Auth, Tracing) có thể bọc lấy nhau và thực thi theo thứ tự.

**Đề bài:** Định nghĩa kiểu dữ liệu `Handler func(string) string`.
Viết hàm `Chain(middlewares ...func(Handler) Handler) func(Handler) Handler`.
Hàm này sẽ nhận vào danh sách các middleware và trả về một hàm cho phép "bọc" một `Handler` gốc. Khi `Handler` cuối cùng được gọi, nó phải chạy qua các middleware theo thứ tự truyền vào.

**Test Cases:**
1.  **Input:** `UpperMiddleware`, `TrimMiddleware` bọc `Handler` gốc trả về " hello ". 
    * **Output:** "HELLO" (Trim trước rồi Upper hoặc ngược lại tùy thứ tự chain).
2.  **Input:** Không truyền middleware nào (`Empty Chain`). 
    * **Output:** Giữ nguyên kết quả của `Handler` gốc.
3.  **Input:** `AuthMiddleware` (nếu string không chứa "admin" thì trả về "403"). 
    * **Output:** Kiểm tra tính chặn đứng (short-circuit) của function.

---

## 🔵 Bài 2: Concurrent Task Observer (Thực tế: Monitoring System)
**Bối cảnh:** Trong các hệ thống phân tán, bạn cần theo dõi thời gian thực thi của các hàm mà không làm thay đổi logic bên trong của chúng.

**Đề bài:** Viết hàm `Observe(fn func(int) int, onStart func(), onFinish func(time.Duration)) func(int) int`.
Hàm trả về một "Wrapped Function". Khi hàm này được gọi:
1.  Gọi `onStart` ngay lập tức.
2.  Thực thi `fn`.
3.  Khi `fn` xong, gọi `onFinish` với tham số là tổng thời gian thực thi.
4.  Trả về kết quả của `fn`.

**Test Cases:**
1.  **Input:** Hàm `fn` ngủ 500ms. 
    * **Output:** `onFinish` nhận vào giá trị $\approx 500ms$.
2.  **Input:** Chạy `Observe` trong một vòng lặp 1000 goroutines. 
    * **Output:** Đảm bảo `onStart` và `onFinish` được gọi đúng 1000 lần (kiểm tra race condition).
3.  **Input:** `fn` là hàm đệ quy (ví dụ tính Fibonacci). 
    * **Output:** Chỉ tính thời gian của lời gọi hàm ngoài cùng.

---

## 🟡 Bài 3: Robust API Gateway (Thực tế: Resiliency Patterns)
**Bối cảnh:** API của bên thứ ba thường xuyên gặp lỗi hoặc panic. Bạn cần một hàm "bảo vệ" để hệ thống của bạn không chết chùm.

**Đề bài:** Viết hàm `WithRetry(fn func() (string, error), attempts int) (string, error)`.
Nếu `fn` trả về lỗi hoặc xảy ra **panic**, hàm phải tự động thực hiện lại tối đa `attempts` lần. Nếu vẫn lỗi sau số lần thử, trả về lỗi cuối cùng. Nếu panic ở lần thử cuối, `WithRetry` phải recover và trả về `error`.

**Test Cases:**
1.  **Input:** `fn` lỗi 2 lần đầu, lần 3 thành công, `attempts = 3`. 
    * **Output:** Thành công ở lần 3.
2.  **Input:** `fn` luôn luôn panic, `attempts = 5`. 
    * **Output:** Recover thành công và trả về error "Task failed after 5 attempts".
3.  **Input:** `attempts = 0` hoặc `1`. 
    * **Output:** Chạy đúng 1 lần duy nhất.

---

## 🔴 Bài 4: Functional Database Query Builder (Thực tế: ORM Design)
**Bối cảnh:** Thay vì nối chuỗi SQL nguy hiểm, các thư viện hiện đại dùng Function để xây dựng câu truy vấn (Query Builder).

**Đề bài:** Tạo struct `Query { Where string, Limit int }`.
Viết hàm `BuildQuery(base Query, options ...func(*Query)) Query`.
Triển khai các "Option functions": `WithLimit(n)`, `AndWhere(cond)`, `ClearWhere()`.

**Test Cases:**
1.  **Input:** `BuildQuery(base, WithLimit(10), AndWhere("id > 5"))`. 
    * **Output:** `Limit: 10, Where: "id > 5"`.
2.  **Input:** `BuildQuery(base, AndWhere("a=1"), AndWhere("b=2"))`. 
    * **Output:** Nối chuỗi logic: `Where: "a=1 AND b=2"`.
3.  **Input:** Truyền nhiều Option mâu thuẫn (vị dụ `WithLimit(10)` rồi `WithLimit(20)`). 
    * **Output:** Option cuối cùng phải được áp dụng (ghi đè).

---

## 🟣 Bài 5: Dynamic Rule Engine (Thực tế: Discount/Promotion System)
**Bối cảnh:** Một trang thương mại điện tử cần áp dụng nhiều mã giảm giá cùng lúc (Ví dụ: Giảm 10% cho khách mới AND Giảm tối đa 50k cho đơn hàng trên 200k).

**Đề bài:** Định nghĩa `Rule func(float64) float64`.
Viết hàm `ComposeRules(mode string, rules ...Rule) Rule`.
* Nếu `mode == "MIN"`: Trả về kết quả nhỏ nhất (có lợi cho chủ shop).
* Nếu `mode == "MAX"`: Trả về kết quả lớn nhất (có lợi cho khách hàng).
* Nếu `mode == "SUM"`: Cộng dồn tất cả các mức giảm giá.

**Test Cases:**
1.  **Input:** Giá gốc 100. Rule1: giảm 10, Rule2: giảm 20. Mode "MIN". 
    * **Output:** 70 (vì 100-20=80 > 100-10=90 là sai, giá sau giảm nhỏ nhất là 70). Lưu ý: Logic "có lợi" cần được định nghĩa rõ là giá sau cùng hay số tiền được giảm. Ở đây hãy lấy **giá sau cùng thấp nhất**.
2.  **Input:** Danh sách Rule rỗng. 
    * **Output:** Trả về giá gốc.
3.  **Input:** Rule trả về giá âm (lỗi logic). 
    * **Output:** Hàm `Compose` phải đảm bảo giá không bao giờ < 0.

---