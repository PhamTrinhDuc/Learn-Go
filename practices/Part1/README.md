### 🟢 Bài 1: Tính toán hình học (Toán tử & Biến)
**Yêu cầu:** Nhập vào chiều dài `a` và chiều rộng `b` của một hình chữ nhật (kiểu `float64`). Tính và in ra chu vi và diện tích.

* **Input:** Hai số thực `a` và `b`.
* **Output:** Chu vi và diện tích, cách nhau bởi dấu cách.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `5.0 3.0` | `16 15` |
| #2 | `2.5 4.0` | `13 10` |

---

### 🟡 Bài 2: Chuyển đổi nhiệt độ (Data Types & Casting)
**Yêu cầu:** Viết chương trình chuyển đổi nhiệt độ từ độ Celsius (`C`) sang Fahrenheit (`F`). Công thức: $F = C \times \frac{9}{5} + 32$.
*Lưu ý:* Hãy cẩn thận với phép chia số nguyên trong Go.

* **Input:** Một số thực `C`.
* **Output:** Giá trị `F` (lấy 2 chữ số thập phân).

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `0` | `32.00` |
| #2 | `37` | `98.60` |
| #3 | `-40` | `-40.00` |

---

### 🔴 Bài 3: Hoán vị không dùng biến tạm (Bitwise Operator)
**Yêu cầu:** Cho hai số nguyên `a` và `b`. Hãy hoán đổi giá trị của chúng mà **không sử dụng** biến thứ ba (biến tạm). Sử dụng toán tử XOR (`^`).

* **Input:** Hai số nguyên `a, b`.
* **Output:** Giá trị của `a` và `b` sau khi hoán đổi.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `10 20` | `20 10` |
| #2 | `123 456` | `456 123` |

---

### 🔵 Bài 4: Kiểm tra kiểu dữ liệu (Zero Values)
**Yêu cầu:** Khai báo 4 biến trong Go mà không gán giá trị đầu kỳ: `int`, `float64`, `string`, `bool`. In ra giá trị mặc định (Zero Value) của chúng trên cùng một dòng.

* **Input:** Không có.
* **Output:** Các giá trị mặc định cách nhau bởi dấu cách.

| Test Case | Expected Output |
| :--- | :--- |
| #1 | `0 0 "" false` | *(Lưu ý: "" đại diện cho chuỗi rỗng)* |

---

### 🟣 Bài 5: Tính tổng chữ số (Modulo & Division)
**Yêu cầu:** Nhập vào một số nguyên dương có 3 chữ số. Tính tổng các chữ số của nó.
*Ví dụ:* `123` -> $1 + 2 + 3 = 6$.

* **Input:** Một số nguyên `n` ($100 \le n \le 999$).
* **Output:** Tổng các chữ số.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `123` | `6` |
| #2 | `999` | `27` |
| #3 | `501` | `6` |
