### 🟢 Bài 15: Hệ thống tính Lãi suất kép (Compound Interest)
**Thực tế:** Mô phỏng kế hoạch tiết kiệm ngân hàng.
**Yêu cầu:** Nhập vào số vốn ban đầu $P$, lãi suất hàng năm $r$ (%), và số năm gửi $n$. In ra số dư cuối cùng của từng năm.
* Sử dụng vòng lặp `for` cơ bản.
* Công thức mỗi năm: $P = P + (P \times r / 100)$.

* **Input:** `P (float64)`, `r (float64)`, `n (int)`.
* **Output:** `Năm {i}: {Số dư}` (làm tròn 2 chữ số).

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `100 10 2` | `Năm 1: 110.00`, `Năm 2: 121.00` |

---

### 🟡 Bài 16: Tìm số bị thiếu trong dãy (Data Integrity)
**Thực tế:** Bạn nhận được một danh sách ID nhân viên từ `1` đến `n`, nhưng có một người quên điểm danh.
**Yêu cầu:** Nhập vào số `n` và sau đó là một chuỗi các số từ `1` đến `n` (nhưng thiếu 1 số). Tìm số đó.
* Sử dụng toán tử XOR hoặc tính tổng để tối ưu.
* Sử dụng `fmt.Scan` trong vòng lặp để nhận dữ liệu liên tục.

* **Input:** `n`, sau đó là `n-1` số nguyên.
* **Output:** Số bị thiếu.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `5` rồi `1 2 4 5` | `3` |
| #2 | `3` rồi `3 1` | `2` |

---

### 🔴 Bài 17: Vẽ biểu đồ Histogram (Data Visualization)
**Thực tế:** Hiển thị tần suất dữ liệu bằng ký tự trên Terminal.
**Yêu cầu:** Nhập vào một danh sách các số nguyên. Với mỗi số `x`, in ra một dòng gồm `x` ký tự `*`.
* Nếu `x < 0`, in `Invalid`.
* Sử dụng vòng lặp lồng nhau (Nested Loops).

* **Input:** Một số `n` (số lượng phần tử), sau đó là `n` số nguyên.
* **Output:** Các dòng dấu `*`.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `3` rồi `5 2 4` | `*****`, `**`, `****` |

---

### 🔵 Bài 18: Ước lượng số Pi - Phương pháp Gregory-Leibniz
**Thực tế:** Tính toán số học trong khoa học máy tính.
**Yêu cầu:** Sử dụng vòng lặp để tính gần đúng số $\pi$ theo công thức:
$$\pi = 4 \times (1 - \frac{1}{3} + \frac{1}{5} - \frac{1}{7} + \frac{1}{9} - ...)$$
* Vòng lặp chạy `n` lần (số hạng càng lớn, độ chính xác càng cao).
* Kết hợp `if` để đổi dấu `+` và `-`.

* **Input:** Số nguyên `n` (số lần lặp).
* **Output:** Giá trị $\pi$ lấy 5 chữ số thập phân.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `10` | `3.04184` |
| #2 | `1000` | `3.14059` |

---

### 🟣 Bài 19: Console Menu & Input Validation (Hệ thống thực)
**Thực tế:** Hầu hết các ứng dụng dòng lệnh (CLI) đều cần một vòng lặp vô hạn chờ lệnh người dùng.
**Yêu cầu:** Viết một chương trình giả lập menu:
1. Nhập số `1`: In "Hello Gopher".
2. Nhập số `2`: Nhập tiếp một số `x`, in bình phương của `x`.
3. Nhập số `3`: Thoát chương trình (`break`).
* Bất kỳ số nào khác: In "Lệnh không hợp lệ, nhập lại".
* Sử dụng `for { ... }` (vòng lặp vô hạn).

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `1` | `Hello Gopher` |
| #2 | `2` rồi `4` | `16` |
| #3 | `3` | `Goodbye!` (Kết thúc) |


### 🟢 Bài 20: Phân tích Thừa số Nguyên tố (Prime Factorization)
**Thực tế:** Ứng dụng trong mã hóa RSA và bảo mật hệ thống.
**Yêu cầu:** Nhập một số nguyên dương $N$. Phân tích $N$ thành tích các thừa số nguyên tố theo định dạng $p_1^{a_1} \times p_2^{a_2} \dots$
* **Tối ưu:** Vòng lặp không nên chạy đến $N$ mà chỉ cần đến $\sqrt{N}$.
* Sử dụng vòng lặp `for` lồng nhau và toán tử `%`.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `60` | `2^2 * 3^1 * 5^1` |
| #2 | `97` | `97^1` (Số nguyên tố) |

---

### 🟡 Bài 21: Nén chuỗi ký tự (Run-Length Encoding)
**Thực tế:** Một thuật toán nén dữ liệu cơ bản dùng trong file hình ảnh (như BMP, TIFF).
**Yêu cầu:** Nhập một chuỗi ký tự lặp lại (ví dụ `AAABBCDDDD`). Xuất ra chuỗi đã nén theo quy tắc: `{Ký tự}{Số lần xuất hiện liên tiếp}`.
* Sử dụng vòng lặp `for range` và so sánh ký tự hiện tại với ký tự trước đó.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `AAABBCDDDD` | `A3B2C1D4` |
| #2 | `GGGGG` | `G5` |

---

### 🔴 Bài 22: Chuyển đổi Cơ số (Base Conversion)
**Thực tế:** Chuyển đổi dữ liệu từ Decimal sang Binary, Hexadecimal trong lập trình hệ thống.
**Yêu cầu:** Nhập một số nguyên dương $N$ (hệ 10) và cơ số $B$ ($2 \le B \le 16$). Hãy in ra giá trị của $N$ ở hệ cơ số $B$.
* **Khó khăn:** Với $B > 10$, phải dùng các ký tự `A, B, C, D, E, F`.
* **Ràng buộc:** Không dùng các hàm có sẵn như `strconv.FormatInt`. Hãy dùng vòng lặp và mảng/chuỗi để lưu kết quả tạm.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `255 16` | `FF` |
| #2 | `13 2` | `1101` |

---

### 🔵 Bài 23: Tìm điểm cân bằng (Equilibrium Point)
**Thực tế:** Xử lý mảng dữ liệu lớn (Big Data) để tìm điểm tối ưu.
**Yêu cầu:** Nhập một dãy số. Tìm vị trí $i$ sao cho tổng các số bên trái $i$ bằng tổng các số bên phải $i$. 
* Nếu có nhiều điểm, lấy điểm đầu tiên. Nếu không có, in `-1`.
* **Yêu cầu tối ưu:** Chỉ dùng **tối đa 2 vòng lặp đơn** (không dùng lồng nhau $O(n^2)$).

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `1 7 3 6 5 6` | `3` (tổng trái = tổng phải = 11) |
| #2 | `1 2 3` | `-1` |

---

### 🟣 Bài 24: Giả lập Hệ thống Rút tiền ATM (Greedy Algorithm)
**Thực tế:** Thuật toán thối tiền lẻ tối ưu số lượng tờ tiền.
**Yêu cầu:** Bạn có các loại mệnh giá: `500k, 200k, 100k, 50k, 20k, 10k`. Nhập số tiền muốn rút $S$ (phải là bội số của 10k).
* Hãy in ra số lượng từng tờ tiền sao cho **tổng số tờ là ít nhất**.
* Sử dụng vòng lặp duyệt qua danh sách mệnh giá.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `880000` | `500k:1, 200k:1, 100k:1, 50k:1, 20k:1, 10k:1` |
| #2 | `2000000` | `500k:4` |
