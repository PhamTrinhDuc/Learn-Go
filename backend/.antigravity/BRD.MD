
# Business Requirements Document
## Hệ thống quản lý chuỗi barber shop tích hợp AI Agent — v2.0

---

### 1. Tổng quan

Hệ thống quản lý chuỗi barber shop với AI Agent là nền tảng cốt lõi. Thay vì AI chỉ là tính năng phụ, toàn bộ trải nghiệm khách hàng được dẫn dắt qua hội thoại với agent — từ tìm hiểu dịch vụ, đặt lịch, nhận nhắc nhở, đến chăm sóc sau dịch vụ. Nhân viên và quản lý có ứng dụng riêng để vận hành.

**Mục tiêu:**
- Tự động hóa toàn bộ hành trình khách hàng qua AI Agent
- Quản lý vận hành chuỗi 5–10 chi nhánh tập trung
- Cung cấp dữ liệu để quản lý ra quyết định

---

### 2. Actors & Phân quyền

| Actor | Mô tả | Quyền hạn |
|---|---|---|
| Khách hàng | Tương tác qua agent các kênh chat | Đặt/hủy/đổi lịch của bản thân, mua hàng, xem lịch sử cá nhân |
| AI Agent | Nhân viên tư vấn tự động | Thực hiện hành động thay mặt khách sau khi có xác nhận HITL |
| Stylist | Nhân viên thực hiện dịch vụ | Xem lịch hẹn được phân công, checkin/checkout khách, xem lịch của bản thân |
| Quản lý chi nhánh | Điều hành một chi nhánh | Mọi quyền của stylist + quản lý ca, xem tất cả booking chi nhánh, xem báo cáo chi nhánh, quản lý KB chi nhánh, duyệt booking đặc biệt |
| Quản lý chuỗi (HQ) | Điều hành toàn hệ thống | Mọi quyền của quản lý chi nhánh + quản lý danh mục dịch vụ/sản phẩm toàn chuỗi, xem báo cáo tổng hợp, quản lý KB toàn chuỗi |

**Lưu ý phân quyền:**
- Stylist chỉ xem được booking được phân công cho mình, không xem booking của stylist khác
- Quản lý chi nhánh chỉ xem được dữ liệu chi nhánh mình quản lý
- HQ xem được tất cả

---

### 3. Yêu cầu nghiệp vụ

---

#### 3.1 Xác thực & tài khoản

**BR01 — Đăng ký**

Khách đăng ký bằng số điện thoại. Luồng:
1. Nhập số điện thoại
2. Hệ thống gửi OTP (6 số, hiệu lực 5 phút)
3. Nhập OTP → tạo tài khoản
4. Nhập thêm tên, ngày sinh (tùy chọn)

Mỗi số điện thoại chỉ đăng ký được một tài khoản.

**BR02 — Đăng nhập**

Đăng nhập bằng số điện thoại + OTP, hoặc số điện thoại + mật khẩu nếu đã thiết lập. Hệ thống cấp JWT token, hiệu lực 7 ngày, refresh tự động khi còn dưới 1 ngày.

**BR03 — Liên kết kênh chat**

Lần đầu khách nhắn tin qua Zalo/Telegram/Discord, agent gửi link xác thực. Khách đăng nhập qua link đó → kênh chat được bind với tài khoản. Từ đó về sau nhắn tin qua kênh này tự động nhận diện không cần đăng nhập lại. Một tài khoản có thể bind nhiều kênh.

---

#### 3.2 Trải nghiệm khách hàng qua AI Agent

**BR04 — Đa kênh hội thoại**

Hỗ trợ: Zalo, Telegram, Discord, Web chat. Cùng một khách nhắn từ bất kỳ kênh nào đã bind đều được nhận diện là một người. Lịch sử hội thoại lưu theo tài khoản, không phải theo kênh — khách bắt đầu trên Zalo, tiếp tục trên Web vẫn thấy context cũ.

**BR05 — Tư vấn dịch vụ & sản phẩm**

Agent tra cứu Knowledge Base để trả lời câu hỏi về dịch vụ, giá, chính sách, địa chỉ chi nhánh. Nếu câu hỏi liên quan đến chi nhánh cụ thể thì ưu tiên KB của chi nhánh đó, fallback sang KB toàn chuỗi.

**BR06 — Đặt lịch qua hội thoại (có HITL)**

Agent dẫn dắt khách qua các bước:
1. Xác định chi nhánh
2. Chọn dịch vụ hoặc combo
3. Chọn ngày
4. Hệ thống trả về danh sách stylist có ca trong ngày đó
5. Khách chọn stylist hoặc chọn "bất kỳ"
6. Hệ thống trả về các slot giờ trống của stylist đó
7. Khách chọn giờ
8. **HITL:** Agent tổng hợp toàn bộ thông tin, hiển thị để khách xác nhận
9. Khách bấm xác nhận → hệ thống tạo booking, trả về mã booking

Nếu chọn "bất kỳ stylist", hệ thống tự phân công stylist có ít booking nhất trong ca đó.

Kiểm tra slot phải theo thời gian thực — không được confirm lịch nếu slot vừa bị người khác book.

**BR07 — Đổi lịch qua hội thoại (có HITL)**

Khách yêu cầu đổi lịch → agent hỏi muốn đổi gì (ngày, giờ, stylist, dịch vụ) → thực hiện tương tự luồng đặt lịch cho thông tin mới → HITL xác nhận → cập nhật. Booking cũ chuyển sang `rescheduled`, booking mới được tạo với `rescheduled_from` trỏ về booking cũ. Không tính vào số lần hủy.

**BR08 — Hủy lịch qua hội thoại (có HITL)**

Khách yêu cầu hủy → agent hỏi lý do (tùy chọn) → **HITL:** "Bạn có chắc muốn hủy booking [mã] vào [ngày giờ] không?" → khách xác nhận → hủy. Ghi nhận `cancel_reason`, `cancel_by = user`.

**Chính sách vi phạm hủy lịch:**

| Tình huống | Hệ quả |
|---|---|
| Hủy lần 1–2 trong 30 ngày | Ghi nhận, không hạn chế |
| Hủy lần 3 trong 30 ngày | Cảnh báo, tài khoản bị hạn chế đặt lịch 7 ngày |
| Hủy lần 4–5 tích lũy | Cảnh báo thêm |
| Hủy lần 5+ tích lũy | Mọi booking mới cần quản lý chi nhánh duyệt thủ công |
| 1 lần no-show | Tính bằng 2 lần hủy |

Vi phạm được đếm theo tài khoản, không theo kênh. Agent thông báo rõ hệ quả trước khi khách xác nhận hủy.

**BR09 — Nhắc lịch tự động**

Hệ thống tạo 2 AGENT_TASK khi booking được confirm:
- Task 1: gửi reminder 24 giờ trước giờ hẹn
- Task 2: gửi reminder 2 giờ trước giờ hẹn

Nội dung reminder gồm: tên khách, tên stylist, dịch vụ, địa chỉ chi nhánh, giờ hẹn, link để xác nhận/đổi/hủy.

Khi booking bị hủy hoặc đổi lịch, các task reminder pending liên quan bị cancel ngay lập tức.

**BR10 — Chăm sóc sau dịch vụ**

Khi booking chuyển sang `completed`, hệ thống tạo AGENT_TASK gửi tin nhắn xin đánh giá sau 1 giờ. Nếu khách không phản hồi thì không nhắc lại.

Nếu khách không có booking nào trong 30 ngày, hệ thống tạo AGENT_TASK gửi tin nhắn tái kích hoạt. Chỉ gửi tối đa 1 lần mỗi 30 ngày, không gửi nếu khách đã tắt marketing.

**BR11 — Escalate sang staff**

Khi gặp: khiếu nại, yêu cầu hoàn tiền, câu hỏi ngoài KB, booking cần duyệt thủ công — agent thông báo khách sẽ được hỗ trợ bởi nhân viên, sau đó tạo notification cho quản lý chi nhánh tương ứng. Quản lý nhận notification trên app, vào xem conversation và tiếp tục xử lý thủ công. Conversation được đánh dấu `escalated`, agent không tự trả lời thêm cho đến khi staff đánh dấu `resolved`.

---

#### 3.3 Quản lý lịch hẹn

**BR12 — Vòng đời booking**

```
pending → confirmed → checked_in → completed
                    ↘ no_show
         ↘ cancelled (từ pending hoặc confirmed)
         ↘ rescheduled (từ pending hoặc confirmed)
```

- `pending`: vừa tạo, chờ hệ thống xử lý (thường chỉ vài giây)
- `confirmed`: đã xác nhận slot, reminder đã được tạo
- `checked_in`: khách đã đến, stylist đang phục vụ
- `completed`: dịch vụ hoàn thành, điểm loyalty được cộng
- `no_show`: quá 15 phút sau giờ hẹn không checkin
- `cancelled`: đã hủy bởi user, staff, hoặc hệ thống
- `rescheduled`: đã đổi sang booking mới

**BR13 — Checkin**

Stylist hoặc quản lý mở app, tìm booking theo mã hoặc tên khách, bấm checkin. Trạng thái chuyển `checked_in`. Chỉ checkin được trong khoảng 15 phút trước đến 30 phút sau giờ hẹn.

**BR14 — Hoàn thành dịch vụ**

Sau khi xong, stylist bấm "Hoàn thành" trên app. Trạng thái chuyển `completed`. Điểm loyalty tự động cộng vào tài khoản khách. Task xin đánh giá được tạo.

**BR15 — Ghi nhận no-show tự động**

Cronjob chạy mỗi 5 phút, tìm các booking `confirmed` có giờ hẹn đã qua 15 phút mà chưa checkin → chuyển sang `no_show`. Cộng điểm vi phạm theo chính sách hủy.

---

#### 3.4 Quản lý sản phẩm & đơn hàng

**BR16 — Danh mục sản phẩm**

Sản phẩm có: tên, danh mục, mô tả, giá niêm yết, giá sale (nullable), ảnh, trạng thái active/inactive. Tồn kho quản lý theo từng chi nhánh.

**BR17 — Đặt hàng (có HITL)**

Khách chọn sản phẩm qua chat hoặc web → agent tổng hợp giỏ hàng → HITL xác nhận → tạo order. Order gắn với chi nhánh gần nhất hoặc khách chỉ định. Thanh toán tại quầy khi nhận hàng (mock — không tích hợp cổng thanh toán thật).

**Vòng đời order:**
```
pending → confirmed → delivering → completed
        ↘ cancelled
```

**BR18 — Tồn kho**

Khi order chuyển `confirmed`, tồn kho tự động giảm. Nếu hủy order, tồn kho hoàn lại. Hệ thống cảnh báo quản lý khi tồn kho xuống dưới `min_stock`. Quản lý nhập hàng thủ công qua app, mỗi lần nhập ghi vào INVENTORY_LOG.

---

#### 3.5 Vận hành chi nhánh

**BR19 — Quản lý ca làm**

Quản lý xếp ca cho stylist theo ngày: ca sáng, ca chiều, hoặc cả ngày. Mỗi ca có `start_time` và `end_time`. Hệ thống chỉ cho phép đặt lịch với stylist trong ca làm của ngày đó, và chỉ trong khung giờ của ca.

**BR20 — Màn hình lịch hẹn**

App hiển thị danh sách booking trong ngày theo timeline. Cập nhật realtime khi có booking mới, hủy, hoặc thay đổi trạng thái. Stylist chỉ thấy booking của mình, quản lý thấy tất cả.

**BR21 — Knowledge Base**

KB gồm các entry: tiêu đề, nội dung, danh mục, embedding vector. Có thể gắn với toàn chuỗi (branch_id = null) hoặc một chi nhánh cụ thể. Khi agent query KB, ưu tiên KB chi nhánh trước, fallback sang KB toàn chuỗi. Quản lý chi nhánh quản lý KB của chi nhánh mình. HQ quản lý KB toàn chuỗi.

---

#### 3.6 Loyalty

**BR22 — Tích điểm**

| Hành động | Điểm tích |
|---|---|
| Booking hoàn thành | 10 điểm |
| Mua sản phẩm | 1 điểm / 10,000 VND |
| Sinh nhật | 50 điểm (cộng vào ngày sinh nhật) |

Điểm không có hạn sử dụng. Mỗi lần cộng/trừ điểm đều ghi vào `LOYALTY_TRANSACTION` với lý do.

**BR23 — Đổi điểm**

| Điểm | Quy đổi |
|---|---|
| 50 điểm | Giảm 10,000 VND cho booking tiếp theo |
| 100 điểm | Giảm 25,000 VND cho booking tiếp theo |
| 200 điểm | 1 dịch vụ cắt tóc cơ bản miễn phí |

Khách thông báo với agent muốn dùng điểm → agent kiểm tra số dư → HITL xác nhận → trừ điểm, ghi nhận vào booking.

---

#### 3.7 Notification & Consent

**BR24 — Phân loại thông báo**

| Loại | Mô tả | Có thể tắt không |
|---|---|---|
| `booking_confirm` | Xác nhận đặt lịch | Không |
| `booking_reminder` | Nhắc lịch | Không |
| `booking_cancelled` | Thông báo hủy | Không |
| `follow_up` | Xin đánh giá sau dịch vụ | Có |
| `re_engage` | Tái kích hoạt | Có |
| `marketing` | Khuyến mãi | Có |
| `loyalty` | Thông báo điểm | Có |

**BR25 — Consent**

Khi đăng ký, mặc định bật tất cả. Khách có thể tắt từng loại thông báo có thể tắt bất kỳ lúc nào qua chat ("tôi không muốn nhận tin khuyến mãi nữa") hoặc qua cài đặt tài khoản. Mọi thay đổi consent ghi vào `USER_CONSENT_LOG`.

**BR26 — Do Not Disturb**

Khách có thể thiết lập giờ không muốn nhận thông báo (ví dụ 22:00 – 07:00). Thông báo rơi vào khung này sẽ được queue và gửi sau khi hết giờ DND. Ngoại lệ: không có loại nào được gửi trong giờ DND kể cả booking_confirm — thay vào đó gửi ngay khi hết DND.

---

#### 3.8 Báo cáo

**BR27 — Báo cáo chi nhánh** (quản lý chi nhánh xem)

- Số booking theo ngày/tuần/tháng, breakdown theo trạng thái
- Tỉ lệ hoàn thành / hủy / no-show
- Hiệu suất từng stylist: số booking hoàn thành, rating trung bình
- Doanh thu từ sản phẩm theo kỳ

**BR28 — Báo cáo toàn chuỗi** (HQ xem)

- Tổng hợp tất cả chỉ số trên theo từng chi nhánh
- So sánh hiệu suất giữa các chi nhánh
- Tổng số khách mới / khách quay lại theo kỳ

---

### 4. Yêu cầu phi chức năng

**Phản hồi Agent:** Câu hỏi thông thường ≤ 3 giây. Tác vụ phức tạp (kiểm tra slot, tạo booking) ≤ 10 giây.

**HITL:** Mọi hành động tác động dữ liệu (tạo/sửa/xóa) đều phải có bước xác nhận tường minh từ người dùng. Agent không được tự thực hiện.

**Realtime:** Màn hình lịch hẹn của staff cập nhật realtime khi có thay đổi.

**Đa kênh:** Thêm kênh chat mới không ảnh hưởng đến core logic của agent.

**Idempotency:** Nếu khách bấm xác nhận 2 lần (do lag), hệ thống chỉ tạo 1 booking.

---

### 5. Ngoài phạm vi

- Thanh toán online thực tế (chỉ mock flow)
- App mobile native (web responsive / PWA)
- Tích hợp máy POS
- Quản lý tài chính, kế toán
- Tuyển dụng, chấm công nhân sự
- Đa ngôn ngữ

---

### 6. Bảng thuật ngữ

| Thuật ngữ | Định nghĩa |
|---|---|
| HITL | Human-in-the-loop — bước xác nhận bắt buộc từ người dùng trước khi agent thực hiện |
| Slot | Khung giờ trống của một stylist trong một ca làm |
| Escalate | Agent chuyển conversation sang staff thật xử lý |
| KB | Knowledge Base — kho thông tin agent dùng để trả lời |
| No-show | Khách đặt lịch nhưng không đến và không hủy trước |
| Rescheduled | Đổi lịch — tạo booking mới thay thế booking cũ |
| DND | Do Not Disturb — khung giờ không nhận thông báo |

---