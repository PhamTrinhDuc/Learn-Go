
## PHẦN 2: FLOW CONTROL (IF/ELSE & SWITCH)

Chào mừng bạn đến với "trái tim" của logic lập trình. Ở phần này, chúng ta sẽ luyện tập cách điều hướng chương trình.

### 🟢 Bài 6: Phân loại tam giác (If/Else)
**Yêu cầu:** Nhập vào 3 số nguyên `a, b, c` là độ dài 3 cạnh.
1. Kiểm tra xem có tạo thành tam giác không ($a+b>c$ và tương tự).
2. Nếu có, phân loại: `Deu` (Đều), `Can` (Cân), hoặc `Thuong` (Thường).
3. Nếu không, in `Invalid`.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `3 3 3` | `Deu` |
| #2 | `3 4 5` | `Thuong` |
| #3 | `1 1 10` | `Invalid` |

---

### 🟡 Bài 7: Tính thuế thu nhập (Nested If)
**Yêu cầu:** Nhập vào thu nhập tháng `m` (triệu VNĐ). Tính số thuế phải nộp theo quy tắc:
* Dưới 11tr: `0%`
* Từ 11tr đến 20tr: `10%` phần vượt quá 11tr.
* Trên 20tr: `900k` (thuế của mức 2) + `20%` phần vượt quá 20tr.

* **Input:** Một số thực `m`.
* **Output:** Số tiền thuế (triệu VNĐ), làm tròn 2 chữ số thập phân.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `10` | `0.00` |
| #2 | `15` | `0.40` | (4tr vượt x 10% = 0.4) |
| #3 | `25` | `1.90` | (0.9 + 5tr vượt x 20% = 1.9) |

---

### 🔴 Bài 8: Máy tính đơn giản (Switch Case)
**Yêu cầu:** Nhập vào hai số thực `a, b` và một ký tự toán tử `op` (`+`, `-`, `*`, `/`). Thực hiện phép tính. Nếu là phép chia, hãy kiểm tra mẫu số bằng 0.

* **Input:** `a (float64)`, `op (string)`, `b (float64)`.
* **Output:** Kết quả phép tính hoặc "Error" nếu chia cho 0 hoặc toán tử không hợp lệ.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `10 + 5` | `15` |
| #2 | `10 / 0` | `Error` |
| #3 | `7 * 3` | `21` |

---

### 🔵 Bài 9: Short Statement If (Đặc sản của Go)
**Yêu cầu:** Go cho phép viết `if v := value; v > limit { ... }`.
Viết chương trình nhập vào một số `x`. Tính $y = x \times 2$. Nếu $y > 100$, in `High`, ngược lại in `Low`. **Bắt buộc** khai báo `y` ngay trong dòng `if`.

* **Input:** Một số nguyên `x`.
* **Output:** `High` hoặc `Low`.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `60` | `High` |
| #2 | `40` | `Low` |


### 🟢 Bài 10: Giải phương trình bậc 2 (Toán học & Logic)
**Yêu cầu:** Giải phương trình $ax^2 + bx + c = 0$.
* Sử dụng gói `math` để tính căn bậc hai (`math.Sqrt`).
* Xét đầy đủ các trường hợp: $a=0$ (trở thành pt bậc 1), $\Delta < 0$ (vô nghiệm), $\Delta = 0$ (nghiệm kép), $\Delta > 0$ (2 nghiệm phân biệt).

* **Input:** 3 số thực `a, b, c`.
* **Output:** * `Vô nghiệm`
    * `Vô số nghiệm`
    * `x = {giá trị}` (nếu 1 nghiệm)
    * `x1 = {val1}, x2 = {val2}` (2 nghiệm, x1 < x2, làm tròn 2 chữ số).

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `0 0 5` | `Vô nghiệm` |
| #2 | `1 -3 2` | `x1 = 1.00, x2 = 2.00` |
| #3 | `1 2 1` | `x = -1.00` |

---

### 🟡 Bài 11: Kiểm tra Năm Nhuận & Ngày trong tháng (Nested Logic)
**Yêu cầu:** Nhập vào `tháng` và `năm`. In ra số ngày của tháng đó.
* **Quy tắc năm nhuận:** Chia hết cho 400 HOẶC (chia hết cho 4 VÀ không chia hết cho 100).
* Sử dụng `switch` cho tháng và `if/else` lồng nhau cho tháng 2.

* **Input:** Hai số nguyên `month, year`.
* **Output:** Số ngày (ví dụ: `28`, `29`, `30`, `31`) hoặc `Invalid` nếu tháng không hợp lệ.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `2 2024` | `29` |
| #2 | `2 2100` | `28` |
| #3 | `11 2023` | `30` |

---

### 🔴 Bài 12: Type Switch (Interface & Types)
**Yêu cầu:** Đây là một tính năng cực mạnh của Go. Giả sử bạn có một biến `i interface{}`. Hãy dùng `switch i.(type)` để kiểm tra kiểu dữ liệu của nó.
* Nếu là `int`: in `Gấp đôi: {i*2}`.
* Nếu là `string`: in `Độ dài: {len(i)}`.
* Các kiểu khác: in `Unknown`.

* **Lưu ý:** Bạn cần ép kiểu (type assertion) để tính toán sau khi xác định được kiểu.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `i = 10` | `Gấp đôi: 20` |
| #2 | `i = "Golang"` | `Độ dài: 6` |

---

### 🔵 Bài 13: Game Oẳn Tù Tì (Random & Logic)
**Yêu cầu:** Viết logic cho một lượt chơi. Người chơi nhập `Kéo`, `Búa`, hoặc `Bao`. Máy tính sẽ chọn ngẫu nhiên.
* Sử dụng `math/rand` để giả lập máy chọn.
* Sử dụng `switch` không có biểu thức (`switch { ... }`) thay vì `if/else` dài dòng để so sánh kết quả.

* **Input:** `string` (Kéo/Búa/Bao).
* **Output:** `Thắng`, `Thua`, hoặc `Hòa` kèm theo lựa chọn của máy.

| Test Case | Input | Expected Output (Ví dụ) |
| :--- | :--- | :--- |
| #1 | `Búa` | `Hòa (Máy chọn Búa)` |

---

### 🟣 Bài 14: Validation Password (Complex Conditions)
**Yêu cầu:** Kiểm tra một chuỗi password có hợp lệ không. Một password mạnh phải thỏa mãn:
1. Độ dài ít nhất 8 ký tự.
2. Có ít nhất 1 chữ hoa và 1 chữ thường.
3. Có ít nhất 1 số.
4. Có ít nhất 1 ký tự đặc biệt trong tập: `@`, `#`, `$`, `%`.

* **Gợi ý:** Sử dụng gói `unicode` (ví dụ: `unicode.IsUpper(r)`) và vòng lặp `for range` để duyệt chuỗi, kết hợp các biến boolean để đánh dấu.

| Test Case | Input | Expected Output |
| :--- | :--- | :--- |
| #1 | `abc123` | `Weak` |
| #2 | `Abc@1234` | `Strong` |
