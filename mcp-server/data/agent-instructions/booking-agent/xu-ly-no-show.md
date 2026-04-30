# Booking Agent — Xử lý No-show

## Định nghĩa no-show

Khách được đánh dấu no-show khi:
- Đã quá 15 phút kể từ giờ hẹn
- Không có check-in tại tiệm
- Không liên hệ báo trước

---

## Quy trình xử lý

### Bước 1 — Phát hiện (tự động)
Hệ thống chạy job kiểm tra mỗi 5 phút. Nếu booking có `scheduled_at < NOW() - 15 phút` và `status = confirmed` và `check_in_at IS NULL` → trigger xử lý no-show.

### Bước 2 — Gửi tin nhắn nhắc
Gửi Zalo cho khách:
> "Chúng tôi nhận thấy bạn chưa đến đúng giờ hẹn lúc [giờ]. Nếu bạn đang trên đường, stylist sẽ cố gắng sắp xếp. Nếu bạn cần đổi lịch, hãy nhắn tin để chúng tôi hỗ trợ."

Chờ 10 phút. Nếu không có phản hồi, chuyển sang bước 3.

### Bước 3 — Cập nhật trạng thái
- Update `status = no_show`
- Ghi `agent_action_log` với `action_type = mark_no_show`
- Notify Slack cho manager chi nhánh

### Bước 4 — Ghi nhận lịch sử
No-show được ghi vào lịch sử khách hàng. Không tự động phạt hay trừ điểm — chỉ ghi nhận để analytics theo dõi.

---

## Chính sách với khách no-show nhiều lần

- **1–2 lần:** Không có hành động đặc biệt, ghi nhận trong lịch sử.
- **3 lần trở lên:** Gửi Slack notify cho owner để quyết định có áp dụng chính sách đặt cọc hay không. Agent không tự quyết.

---

## Lưu ý

Agent không được tự huỷ lịch của stylist mà không có xác nhận từ manager hoặc owner. Stylist vẫn được tính ca làm việc bình thường dù khách no-show.
