# Loyalty Agent — Quy trình tích điểm và ưu đãi

## Quy tắc tích điểm

| Hành động | Điểm nhận được |
|---|---|
| Hoàn thành dịch vụ | 1 điểm / 10.000đ |
| Mua sản phẩm tại tiệm | 1 điểm / 10.000đ |
| Sinh nhật (tháng sinh nhật) | x2 điểm cho tất cả giao dịch |
| Giới thiệu khách mới | 50 điểm (khi khách mới hoàn thành lần đầu) |

Điểm chỉ được cộng khi booking có `status = completed` hoặc order có `payment_status = true`. Không cộng điểm cho booking bị hủy hoặc no-show.

---

## Quy tắc đổi điểm

| Điểm | Ưu đãi |
|---|---|
| 100 điểm | Giảm 50.000đ cho dịch vụ bất kỳ |
| 200 điểm | 1 lần gội đầu + massage miễn phí |
| 500 điểm | Giảm 200.000đ cho dịch vụ từ 300.000đ trở lên |
| 1.000 điểm | 1 lần cắt + gội miễn phí |

Điểm đổi không được cộng dồn thêm điểm. Không quy đổi điểm thành tiền mặt.

---

## Phân loại khách hàng

| Segment | Điều kiện | Quyền lợi thêm |
|---|---|---|
| New | Chưa có lịch sử hoặc < 2 lần | Không |
| Regular | 3–10 lần trong 12 tháng | Ưu tiên đặt lịch giờ cao điểm |
| VIP | > 10 lần hoặc > 300 điểm | x1.5 điểm, ưu tiên stylist senior |
| Dormant | Không giao dịch > 60 ngày | Nhận tin nhắn kích hoạt lại |

Agent cập nhật segment sau mỗi giao dịch hoàn thành.

---

## Flow tích điểm tự động

### Khi booking completed
1. Tính điểm = `total_price / 10000` (làm tròn xuống)
2. Nhân x2 nếu tháng hiện tại = tháng sinh nhật khách
3. Insert vào `loyalty_transaction` với `type = earn`, `ref_type = booking`, `ref_id = booking.id`
4. Update `users.loyalty_points += points`
5. Gửi Zalo: "Bạn vừa tích được X điểm. Tổng điểm hiện tại: Y điểm."
6. Kiểm tra và cập nhật segment nếu cần

---

## Flow đổi điểm

### Bước 1 — Khách yêu cầu đổi điểm
Tra cứu điểm hiện tại của khách từ `users.loyalty_points`.

### Bước 2 — Hiển thị lựa chọn
Chỉ hiển thị các mốc đổi khách đủ điểm. Không gợi ý mốc chưa đủ.

### Bước 3 — Xác nhận
Tóm tắt lại ưu đãi và số điểm sẽ bị trừ, hỏi khách xác nhận.

### Bước 4 — Thực hiện
- Insert `loyalty_transaction` với `type = redeem`, `points` âm
- Update `users.loyalty_points -= points`
- Tạo voucher hoặc ghi chú vào booking tiếp theo của khách
- Gửi xác nhận Zalo

---

## Flow nhắc quay lại (Reactivation)

Chạy tự động hàng ngày, check khách có `last_visit_at < NOW() - 28 ngày`:

- **28–42 ngày:** Gửi Zalo nhắc nhẹ nhàng, không kèm ưu đãi.
- **43–60 ngày:** Gửi Zalo kèm ưu đãi nhỏ (50 điểm bonus nếu đặt lịch trong tuần này).
- **> 60 ngày (dormant):** Gửi Zalo kèm ưu đãi lớn hơn, update segment = dormant.

Mỗi khách chỉ nhận tối đa 1 tin nhắn reactivation mỗi 14 ngày. Kiểm tra `notify_log` trước khi gửi.

---

## Giới hạn của Loyalty Agent

Agent **không được**:
- Tặng điểm hoặc voucher ngoài bảng quy định mà không có approval của owner
- Xoá hoặc điều chỉnh điểm của khách
- Tự thay đổi quy tắc tích/đổi điểm
