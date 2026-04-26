## Thiết kế tổng thể hệ thống: Chuỗi cắt tóc multi-agent

---

### Bối cảnh & ràng buộc

Doanh nghiệp 1 người điều hành, nhiều chi nhánh, có nhân viên (stylist) tại mỗi chi nhánh. Owner không muốn xử lý vận hành thủ công hàng ngày — agent lo phần lớn, owner chỉ quyết định những gì quan trọng.

---

## Các actor trong hệ thống

**Khách hàng** — tương tác qua Zalo / Web widget. Đặt lịch, hỏi dịch vụ, mua sản phẩm, nhận thông báo.

**Owner** — tương tác qua Dashboard web + Slack. Phê duyệt action nhạy cảm, xem báo cáo, cấu hình hệ thống.

**Manager chi nhánh** — tương tác qua Mobile app. Cập nhật tồn kho, xem lịch stylist, xử lý vấn đề tại chỗ.

**Stylist** — không tương tác trực tiếp với hệ thống, chỉ nhận lịch làm việc qua app đơn giản.

---

## Module 1 — Booking

### Mục tiêu
Khách tự đặt, đổi, hủy lịch mà không cần gọi điện. Stylist biết lịch của mình. Hệ thống tự nhắc.

### Use cases

**UC-B1: Khách đặt lịch mới**
Khách chọn chi nhánh → chọn dịch vụ → chọn stylist (hoặc để hệ thống chọn) → chọn khung giờ trống → xác nhận. Agent kiểm tra slot thực tế, tạo booking, gửi xác nhận qua Zalo.

**UC-B2: Khách đổi / hủy lịch**
Khách nhắn tin yêu cầu. Agent tìm booking theo số điện thoại, xác nhận thông tin, thực hiện thay đổi. Nếu hủy trong vòng 2 giờ trước giờ hẹn thì gửi Slack notify cho manager chi nhánh.

**UC-B3: Reminder tự động**
Hệ thống tự gửi Zalo nhắc lịch trước 24 giờ và trước 1 giờ. Không cần người kích hoạt.

**UC-B4: Khách không đến (no-show)**
Sau 15 phút quá giờ hẹn mà không có check-in, hệ thống đánh dấu no-show, cộng vào lịch sử khách hàng.

### Quy tắc nghiệp vụ
- Mỗi stylist có lịch làm việc riêng, slot chỉ hiện khi stylist đang làm ca đó
- Một slot = thời lượng dịch vụ + 5 phút buffer
- Không cho đặt lịch trước quá 30 ngày
- Hủy lịch hàng loạt (ví dụ stylist nghỉ đột xuất) → phải qua approve của owner

### Ngoài scope
- Đặt lịch cho nhóm nhiều người
- Thanh toán online khi đặt lịch
- Quản lý ca làm việc của stylist (đây là việc của manager)

---

## Module 2 — FAQ & Tư vấn

### Mục tiêu
Trả lời tự động 80% câu hỏi thường gặp, upsell tự nhiên, không cần nhân viên trực chat.

### Use cases

**UC-F1: Hỏi dịch vụ & giá**
Khách hỏi giá cắt, nhuộm, uốn, v.v. Agent tra Vector Store trả lời đúng bảng giá của chi nhánh đó (vì giá có thể khác nhau giữa chi nhánh).

**UC-F2: Hỏi địa chỉ, giờ mở cửa**
Agent trả lời theo thông tin từng chi nhánh.

**UC-F3: Upsell sau khi tư vấn**
Khách hỏi về dịch vụ cắt → agent gợi ý thêm dịch vụ liên quan (ví dụ: gội đầu + massage) nếu phù hợp với context. Không spam upsell mọi lúc.

**UC-F4: Câu hỏi không có trong knowledge base**
Agent trả lời không biết, hỏi lại khách có muốn để lại số để manager liên hệ không.

### Quy tắc nghiệp vụ
- Knowledge base được owner/manager cập nhật qua Dashboard, không cần kỹ thuật
- Giá và dịch vụ có thể khác nhau theo chi nhánh — agent phải biết đang nói chuyện với khách của chi nhánh nào
- Agent không được bịa thông tin khi không chắc

### Ngoài scope
- Tư vấn kỹ thuật chuyên sâu về tóc (để stylist làm trực tiếp)
- Xử lý khiếu nại — escalate lên owner

---

## Module 3 — Inventory & Bán sản phẩm

### Mục tiêu
Quản lý tồn kho 2 luồng (nội bộ + bán lẻ), cho phép agent bán sản phẩm cho khách, manager kiểm soát kho thực tế.

### Use cases

**UC-I1: Manager nhập hàng / cập nhật tồn kho**
Manager dùng mobile app nhập số lượng nhập mới hoặc điều chỉnh sau khi kiểm kho. Ghi log mọi thay đổi kèm thời gian và người thực hiện.

**UC-I2: Khách mua sản phẩm**
Khách hỏi mua sản phẩm → agent kiểm tra tồn kho bán lẻ → báo giá → khách xác nhận → agent tạo đơn, trừ tồn kho bán lẻ, chuyển thông tin thanh toán.

**UC-I3: Cảnh báo tồn kho thấp**
Khi tồn kho bán lẻ xuống dưới ngưỡng → Slack notify manager chi nhánh. Khi tồn kho nội bộ xuống dưới ngưỡng → Slack notify manager để đặt hàng nhập thêm.

**UC-I4: Tự động tắt bán khi sắp hết**
Khi tồn kho bán lẻ còn dưới ngưỡng tối thiểu (do owner cài đặt) → agent tự ngừng nhận đơn sản phẩm đó, báo khách hàng hết hàng.

### Phân quyền

| Hành động | Owner | Manager | Agent |
|---|---|---|---|
| Nhập hàng mới | Có | Có | Không |
| Điều chỉnh tồn kho nội bộ | Có | Có | Không |
| Bán lẻ / trừ tồn kho | Không | Không | Có |
| Đặt ngưỡng cảnh báo | Có | Không | Không |
| Thay đổi giá bán | Có | Không | Không |
| Đánh dấu sản phẩm được bán lẻ | Có | Không | Không |

### Logic tồn kho
```
Tồn kho thực = Tổng nhập - Bán lẻ (agent) - Nội bộ (manager cập nhật)
```
Tồn kho nội bộ do manager tự trừ định kỳ (theo ca hoặc theo tuần), không tự động hóa vì không có thiết bị đo lường tại chỗ.

### Ngoài scope
- Tự động đặt hàng với nhà cung cấp
- Đồng bộ tồn kho real-time giữa chi nhánh
- Quản lý hạn sử dụng sản phẩm

---

## Module 4 — Loyalty & Giữ chân khách

### Mục tiêu
Tăng tần suất quay lại, tăng giá trị trung bình mỗi khách, giảm chi phí marketing.

### Use cases

**UC-L1: Tích điểm sau mỗi lần dùng dịch vụ / mua hàng**
Mỗi đơn hoàn thành tự động cộng điểm vào tài khoản khách. Khách nhận thông báo Zalo.

**UC-L2: Đổi điểm lấy ưu đãi**
Khách nhắn hỏi điểm của mình, agent tra cứu và hướng dẫn đổi điểm lấy dịch vụ hoặc voucher giảm giá.

**UC-L3: Nhắc quay lại**
Sau 4 tuần kể từ lần cắt tóc gần nhất, hệ thống tự gửi Zalo gợi ý đặt lịch. Nếu khách không phản hồi sau 2 tuần thêm, gửi thêm 1 lần với ưu đãi nhỏ.

**UC-L4: Chăm sóc sinh nhật**
Tự động gửi tin sinh nhật kèm voucher vào ngày sinh nhật khách.

**UC-L5: Phân loại khách**
Hệ thống tự xếp khách vào nhóm (mới, thường xuyên, VIP, ngủ đông) dựa trên lịch sử. Agent điều chỉnh cách giao tiếp theo từng nhóm.

### Quy tắc nghiệp vụ
- Điểm chỉ tích khi đơn được đánh dấu hoàn thành, không tích khi đặt lịch
- Quy tắc tích điểm và đổi điểm do owner cài đặt, không hardcode
- Khách có thể opt-out nhận tin nhắn tự động

### Ngoài scope
- Chương trình giới thiệu bạn bè
- Tích điểm liên chi nhánh kiểu tập đoàn lớn (giữ đơn giản: điểm dùng được ở mọi chi nhánh)

---

## Module 5 — Analytics & Báo cáo

### Mục tiêu
Owner nắm được tình hình toàn chuỗi mà không cần hỏi từng manager.

### Use cases

**UC-A1: Báo cáo doanh thu**
Theo ngày / tuần / tháng, theo chi nhánh, theo dịch vụ, theo stylist. Owner xem trên Dashboard hoặc hỏi agent bằng ngôn ngữ tự nhiên ("doanh thu tuần này chi nhánh Q1 bao nhiêu?").

**UC-A2: Phân tích rush hour**
Xác định khung giờ cao điểm theo chi nhánh để hỗ trợ quyết định xếp ca nhân viên.

**UC-A3: Hiệu suất stylist**
Số lịch hoàn thành, no-show rate, doanh thu mang về. Không ranking công khai — chỉ owner xem.

**UC-A4: Cảnh báo bất thường**
Doanh thu giảm đột ngột so với tuần trước → Slack notify owner tự động.

### Ngoài scope
- Dự báo doanh thu tương lai
- So sánh với đối thủ cạnh tranh
- Báo cáo thuế / kế toán

---

## Human-in-the-loop — Những gì PHẢI qua owner approve

Các action sau agent không được tự thực hiện, phải gửi Slack và chờ:

- Hủy lịch hàng loạt (hơn 3 lịch cùng lúc)
- Hoàn tiền cho khách
- Tặng điểm / voucher ngoài chương trình có sẵn
- Gửi tin nhắn blast cho toàn bộ khách hàng

---

## Kiến trúc dữ liệu — Các entity chính

**Customer**: id, tên, SĐT, ngày sinh, chi nhánh hay đến, điểm tích lũy, nhóm khách, lịch sử tương tác

**Branch**: id, tên, địa chỉ, giờ mở cửa, danh sách stylist, bảng giá riêng

**Booking**: id, customer, branch, stylist, dịch vụ, thời gian, trạng thái, ghi chú

**Product**: id, tên, loại (nội bộ / bán lẻ / cả hai), giá bán, ngưỡng cảnh báo

**Inventory**: product, branch, qty_total, qty_retail, qty_internal, cập nhật lần cuối bởi ai

**Order**: id, customer, items, tổng tiền, trạng thái thanh toán, điểm tích được

**Stylist**: id, tên, branch, lịch làm việc, speciality

---

## Tóm tắt agent nào làm gì

| Agent | Đọc | Ghi | Notify |
|---|---|---|---|
| Booking Agent | Lịch stylist, slot trống | Tạo/sửa/hủy booking | Zalo khách, Slack manager (hủy hàng loạt) |
| FAQ Agent | Knowledge base, bảng giá | Không | Slack nếu câu hỏi không trả lời được |
| Inventory Agent | Tồn kho bán lẻ | Trừ tồn kho khi bán | Slack manager khi thấp |
| Loyalty Agent | Điểm, lịch sử khách | Cộng điểm, tạo voucher | Zalo khách (nhắc lịch, sinh nhật) |
| Analytics Agent | Toàn bộ dữ liệu | Không | Slack owner khi bất thường |