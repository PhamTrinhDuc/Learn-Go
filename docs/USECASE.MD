Dưới đây là danh sách **các Use Case** chính của hệ thống, được chia rõ ràng thành hai phần:

- **Use Case của Website** (Customer-facing): Những chức năng khách hàng trực tiếp tương tác qua web (giao diện chat + form đặt lịch).
- **Use Case của Hệ thống Multi-Agent AI**: Những chức năng mà các Agent (Orchestrator, Booking, RAG/FAQ, Inventory, Analytics…) tự động xử lý ở backend, hỗ trợ chủ tiệm và khách hàng 24/7.

Thiết kế dựa sát trên nghiệp vụ thực tế của **30Shine**: đặt lịch nhanh (chỉ 30 giây), chọn chi nhánh, chọn dịch vụ/combo, hủy lịch thoải mái (không phí), cắt xong trả tiền, quy trình phục vụ cố định, và hỗ trợ multi-chi nhánh.

### 1. Use Case của Website (Khách hàng sử dụng)

| STT | Use Case | Mô tả chi tiết | Actor |
|-----|----------|----------------|-------|
| 1 | Xem danh sách chi nhánh | Xem bản đồ hoặc danh sách chi nhánh gần nhất, địa chỉ, giờ mở cửa | Khách hàng |
| 2 | Xem danh sách dịch vụ & Combo | Xem tất cả dịch vụ (cắt tóc, gội dưỡng, lấy ráy tai…) và các combo đang hot (Shine Combo, Giữ nếp X2, Bay mụn, Dưỡng da…) kèm giá, thời lượng, mô tả | Khách hàng |
| 3 | Đặt lịch hẹn mới | Chọn chi nhánh → Chọn dịch vụ/combo → Chọn ngày/giờ trống → Nhập tên + số điện thoại → Xác nhận đặt lịch (chỉ cần phone là chính) | Khách hàng |
| 4 | Chọn stylist (tùy chọn) | Chọn nhân viên thực hiện dịch vụ (nếu muốn) | Khách hàng |
| 5 | Xem lịch hẹn của mình | Tra cứu lịch đã đặt theo số điện thoại | Khách hàng |
| 6 | Hủy lịch hẹn | Hủy lịch đã đặt (hệ thống cho phép hủy thoải mái, đặc biệt nếu chưa gần giờ) | Khách hàng |
| 7 | Hỏi thông tin dịch vụ & giá | Chat hỏi giá combo, chi tiết dịch vụ, quy trình phục vụ (7-10 bước) | Khách hàng |
| 8 | Hỏi chính sách tiệm | Hỏi giờ mở cửa, chính sách hủy lịch, thanh toán (cắt xong trả tiền), chính sách đến muộn | Khách hàng |
| 9 | Hỏi thông tin chung | Hỏi địa chỉ chi nhánh, thời gian phục vụ tốt nhất, dịch vụ mới… | Khách hàng |
| 10 | Nhận xác nhận & nhắc lịch | Nhận thông báo đặt lịch thành công + nhắc lịch qua website/chat (hoặc SMS nếu tích hợp sau) | Khách hàng |
| 11 | Đánh giá dịch vụ (tùy chọn) | Sau khi hoàn thành, mời đánh giá chất lượng phục vụ | Khách hàng |

**Lưu ý**: Website ưu tiên **flow cực kỳ đơn giản** như 30Shine (đặt lịch nhanh, không bắt buộc tạo tài khoản đầy đủ).

### 2. Use Case của Hệ thống Multi-Agent AI

| STT | Use Case | Agent chính chịu trách nhiệm | Mô tả chi tiết |
|-----|----------|------------------------------|--------------|
| 1 | Nhận và phân loại intent từ khách | Orchestrator Agent | Nhận tin nhắn từ website/chat, phân loại intent (đặt lịch / hỏi giá / hủy lịch / hỏi chính sách…) bằng rule-based + LLM |
| 2 | Xử lý đặt lịch hẹn mới | Booking Agent + Orchestrator | Check slot trống theo chi nhánh + stylist + thời gian → Tạo booking (pending) → Gửi soft confirm qua Slack cho chủ → Auto confirm nếu không reject |
| 3 | Trả lời câu hỏi FAQ & thông tin dịch vụ | FAQ/RAG Agent | Sử dụng RAG để trả lời giá, combo, quy trình phục vụ, chính sách dựa trên knowledge base |
| 4 | Xử lý hủy lịch | Booking Agent | Kiểm tra lịch → Thực hiện hủy (auto nếu xa giờ, cần approval nếu gần giờ) |
| 5 | Nhắc lịch cho khách | Booking Agent (Scheduler) | Tự động gửi nhắc lịch trước 1 giờ (qua website hoặc kênh khác) |
| 6 | Cảnh báo tồn kho thấp | Inventory Agent | Theo dõi vật tư (dầu gội, sáp, khăn…) → Gửi alert qua Slack cho chủ khi sắp hết |
| 7 | Trừ vật tư sau khi hoàn thành dịch vụ | Inventory Agent | Khi booking chuyển sang **completed**, tự động ghi log và trừ tồn kho dựa trên lượng tiêu hao trung bình |
| 8 | Gửi báo cáo hàng ngày/tuần/tháng | Analytics Agent | Gửi qua Slack: số lịch hoàn thành, doanh thu ước tính, no-show rate, dịch vụ bán chạy, tồn kho cần nhập… |
| 9 | Quản lý approval từ chủ tiệm | Orchestrator Agent | Gửi thông báo Slack với nút Approve/Reject cho các action có rủi ro (tạo lịch mới, hủy gần giờ, thay đổi lớn) |
| 10 | Quản lý memory & lịch sử hội thoại | Orchestrator + Conversation table | Lưu lịch sử chat ngắn hạn để duy trì ngữ cảnh cuộc trò chuyện |
| 11 | Xử lý trường hợp đặc biệt (fallback) | Orchestrator | Khi RAG không tìm được thông tin → Hỏi chủ tiệm qua Slack hoặc gợi ý đặt lịch |
| 12 | Rate limiting & chống spam | Orchestrator | Giới hạn số lần đặt/hỏi theo số điện thoại để tránh lạm dụng |

### Phân loại theo Risk Level (Approval Flow)
- **Auto-run (🟢)**: Trả lời FAQ, nhắc lịch, báo cáo, trừ kho sau completed.
- **Soft confirm (🟡)**: Tạo lịch hẹn mới (notify Slack, auto confirm sau 5-10 phút nếu chủ không phản hồi).
- **Hard approval (🔴)**: Hủy lịch gần giờ, thay đổi thông tin lớn.

Bạn muốn tôi mở rộng chi tiết bất kỳ use case nào (ví dụ: sequence diagram cho “Đặt lịch hẹn mới” hoặc “Hỏi giá combo”), hay liệt kê thêm use case cho chủ tiệm qua Slack?

Hoặc bạn cần sắp xếp theo mức độ ưu tiên cho MVP không?