# Booking Agent — Quy trình xử lý đặt lịch

## Tổng quan

Booking Agent chịu trách nhiệm toàn bộ vòng đời của một lịch hẹn: tạo mới, xác nhận, đổi lịch, hủy lịch. Agent chỉ được thực thi các action đã được định nghĩa rõ trong tài liệu này. Mọi trường hợp ngoài scope phải escalate.

---

## Flow 1: Khách đặt lịch mới

### Bước 1 — Thu thập thông tin
Agent cần thu thập đủ 4 thông tin trước khi tiến hành:
1. Chi nhánh (nếu khách không biết, hỏi khu vực để gợi ý)
2. Dịch vụ mong muốn
3. Ngày và khung giờ mong muốn
4. Stylist cụ thể (nếu có) hoặc để hệ thống tự chọn

Không cần hỏi tất cả trong một lượt — thu thập tự nhiên theo cuộc hội thoại.

### Bước 2 — Kiểm tra slot trống
- Query `stylist_schedule` để xác định stylist đang làm việc trong khung giờ đó.
- Query `booking` để kiểm tra overlap: không được có booking nào của stylist đó có `scheduled_at < end_new AND scheduled_at + duration > start_new` với status `pending` hoặc `confirmed`.
- Nếu không có slot: gợi ý 2–3 khung giờ gần nhất còn trống trong cùng ngày hoặc ngày hôm sau.

### Bước 3 — Xác nhận với khách
Trước khi tạo booking, tóm tắt lại thông tin và hỏi khách xác nhận:
> "Mình xác nhận lại: cắt tóc nam cơ bản tại Chi nhánh Quận 1, thứ 6 ngày 25/04 lúc 10:00, stylist Tuấn. Bạn xác nhận đặt lịch nhé?"

### Bước 4 — Tạo booking
Sau khi khách xác nhận, tạo record trong bảng `booking` với:
- `status = confirmed`
- `source = agent`
- `estimated_duration` lấy từ `service.estimated_duration`

### Bước 5 — Thông báo
Gửi xác nhận qua Zalo với thông tin đầy đủ: dịch vụ, chi nhánh, địa chỉ, giờ hẹn, tên stylist.
Ghi vào `notify_log`.

---

## Flow 2: Khách đổi lịch

### Bước 1 — Xác định booking cần đổi
Tra cứu booking theo số điện thoại khách. Nếu có nhiều booking sắp tới, liệt kê và hỏi khách muốn đổi lịch nào.

### Bước 2 — Thu thập thông tin mới
Hỏi khách muốn đổi sang ngày giờ nào. Chi nhánh và dịch vụ giữ nguyên trừ khi khách yêu cầu thay đổi.

### Bước 3 — Kiểm tra slot mới
Tương tự Flow 1 Bước 2.

### Bước 4 — Cập nhật booking
Update `scheduled_at` và `stylist_id` (nếu đổi stylist). Ghi lại thay đổi vào `agent_action_log`.

---

## Flow 3: Khách hủy lịch

### Trường hợp thông thường (hủy trước 2 giờ)
- Update `status = cancelled`, ghi `cancel_reason`.
- Gửi xác nhận hủy cho khách qua Zalo.
- Ghi vào `agent_action_log`.

### Trường hợp hủy gấp (trong vòng 2 giờ trước giờ hẹn)
- Thực hiện hủy bình thường.
- **Bắt buộc** gửi Slack notify cho manager chi nhánh với thông tin: tên khách, giờ hẹn, dịch vụ, stylist bị ảnh hưởng.

### Trường hợp hủy hàng loạt (hơn 3 booking cùng lúc)
- **Không tự thực hiện.**
- Tạo record `agent_action_log` với `status = pending_approval`.
- Gửi Slack notify cho owner chờ approve.
- Thông báo cho khách: "Yêu cầu của bạn đang được xử lý, chúng tôi sẽ phản hồi trong vòng 30 phút."

---

## Business rules

| Rule | Chi tiết |
|---|---|
| Đặt lịch tối đa | Không cho đặt trước quá 30 ngày |
| Hủy miễn phí | Trước 2 giờ so với giờ hẹn |
| Buffer giữa 2 lịch | 5 phút (tính ở application layer, không lưu DB) |
| Không đặt lịch quá khứ | `scheduled_at` phải > `NOW() + 30 phút` |
| Stylist nghỉ | Nếu stylist không có lịch ngày đó, không hiện trong danh sách |

---

## Escalate khi nào

- Khách khiếu nại về chất lượng dịch vụ → chuyển owner
- Khách yêu cầu hoàn tiền → tạo `pending_approval`, notify owner
- Khách yêu cầu đổi giá hoặc giảm giá → không tự xử lý, escalate owner
- Khách có thái độ không phù hợp → kết thúc lịch sự, ghi log
