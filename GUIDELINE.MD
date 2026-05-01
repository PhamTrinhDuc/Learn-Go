# Guideline Báo cáo: Hệ thống Agentic AI cho Chuỗi Cắt Tóc

---

## Cấu trúc tổng thể

**Tên đề tài gợi ý:** *Ứng dụng Hệ thống Multi-Agent AI trong Quản lý Vận hành Chuỗi Cửa hàng Cắt tóc cho Doanh nghiệp Một Người*

---

## Chương 1 — Tổng quan bài toán & bối cảnh thực tế

Mục tiêu chương này là trả lời câu hỏi: *"Vấn đề đang tồn tại ở đâu trong thực tế, tại sao nó đau, và tại sao AI là hướng đi hợp lý?"*

**1.1 Bức tranh ngành dịch vụ tóc tại Việt Nam**
Quy mô thị trường, xu hướng chuỗi hóa, áp lực cạnh tranh. Dẫn chứng: mô hình 30Shine, Toni&Guy,... cho thấy xu hướng vận hành bài bản ngay cả ở quy mô nhỏ.

**1.2 Chân dung doanh nghiệp 1 người điều hành nhiều chi nhánh**
Mô tả cụ thể: owner phải đồng thời là giám đốc, kế toán, marketing, HR. Trình bày "một ngày làm việc điển hình" của owner để làm rõ sự quá tải — đây là phần quan trọng nhất của chương, nên dùng hình minh họa hoặc timeline.

**1.3 Các điểm đau cụ thể (Pain Points)**
Trình bày theo 5 nhóm tương ứng với 5 module sau này:

| Nhóm vấn đề | Biểu hiện thực tế | Hậu quả đo lường được |
|---|---|---|
| Đặt lịch thủ công | Gọi điện, nhắn Zalo, ghi sổ tay | Nhầm lịch, mất khách |
| Trực chat 24/7 | Nhân viên trả lời chậm ngoài giờ | Khách bỏ sang đối thủ |
| Quản lý kho mờ nhạt | Không biết tồn kho thực tế | Hết hàng mà không hay |
| Không có loyalty system | Không nhớ khách, không nhắc quay lại | Tần suất quay lại thấp |
| Báo cáo thủ công | Tổng hợp Excel cuối tháng | Quyết định chậm, thiếu dữ liệu |

**1.4 Giới hạn của các giải pháp hiện tại**
Phân tích tại sao phần mềm quản lý salon thông thường (Salonist, KiotViet,...) chưa đủ: không có AI conversational, không tự động hóa workflow, không có human-in-the-loop thông minh.

**1.5 Câu hỏi nghiên cứu & phạm vi đề tài**
Phát biểu rõ bài toán: *"Thiết kế hệ thống multi-agent AI có thể xử lý vận hành hàng ngày của chuỗi cắt tóc, để owner chỉ can thiệp vào các quyết định mang tính chiến lược."*

---

## Chương 2 — Cơ sở lý thuyết & Công nghệ nền tảng

**2.1 Agentic AI & Multi-Agent Systems**
Định nghĩa agent, sự khác biệt giữa AI trả lời đơn thuần vs. AI có khả năng hành động (tool use, memory, planning). Giới thiệu khái niệm orchestrator + sub-agent.

**2.2 Human-in-the-Loop (HITL)**
Tại sao không nên để AI tự quyết hoàn toàn — trình bày triết lý thiết kế: agent xử lý 80%, human quyết định 20% quan trọng.

**2.3 RAG (Retrieval-Augmented Generation)**
Giải thích ngắn tại sao dùng Vector Store cho FAQ thay vì fine-tune model — phù hợp với doanh nghiệp nhỏ cần cập nhật knowledge base linh hoạt.

**2.4 Các công nghệ & nền tảng sử dụng**
Liệt kê và giải thích lý do chọn: LLM (Claude/GPT), Zalo OA API, Slack API, database (PostgreSQL hoặc tương đương), vector database (Pinecone/Qdrant).

---

## Chương 3 — Phân tích & Thiết kế Hệ thống

**3.1 Kiến trúc tổng thể**
Vẽ sơ đồ kiến trúc toàn hệ thống: các actor (Khách, Owner, Manager, Stylist) → các kênh tương tác → Orchestrator → 5 sub-agent → các hệ thống lưu trữ. Đây là hình quan trọng nhất của toàn bộ báo cáo.

**3.2 Thiết kế từng Agent**
Với mỗi trong 5 agent (Booking, FAQ, Inventory, Loyalty, Analytics), trình bày theo template thống nhất:
- Mục tiêu & phạm vi trách nhiệm
- Danh sách use cases (UC)
- Input — Output — Tool sử dụng
- Quy tắc nghiệp vụ quan trọng
- Giới hạn (ngoài scope)

**3.3 Luồng Human-in-the-Loop**
Vẽ flowchart mô tả rõ: agent detect action nhạy cảm → gửi Slack → owner approve/reject → agent thực thi. Các action bắt buộc qua approve: hủy lịch hàng loạt, hoàn tiền, blast message...

**3.4 Thiết kế dữ liệu**
ERD (Entity-Relationship Diagram) cho 7 entity chính: Customer, Branch, Booking, Product, Inventory, Order, Stylist. Chú thích rõ quan hệ và các trường quan trọng.

**3.5 Phân quyền hệ thống**
Bảng phân quyền chi tiết theo từng actor × từng hành động (như bảng trong README của bạn, mở rộng cho toàn hệ thống).

---

## Chương 4 — Kế hoạch Triển khai

**4.1 Phân chia giai đoạn (Roadmap)**

| Giai đoạn | Thời gian | Deliverable |
|---|---|---|
| MVP | Tuần 1–4 | Booking Agent + FAQ Agent hoạt động trên Zalo |
| Beta | Tuần 5–8 | Thêm Inventory + Loyalty, thử nghiệm 1 chi nhánh |
| Full | Tuần 9–12 | Analytics + Dashboard, rollout toàn chuỗi |

**4.2 Yêu cầu hạ tầng & chi phí ước tính**
Ước tính chi phí vận hành hàng tháng: LLM API calls, hosting, Zalo OA. So sánh với chi phí nhân sự tương đương (1 nhân viên trực chat toàn thời gian).

**4.3 Rủi ro & Phương án dự phòng**
Ít nhất 4 rủi ro: AI trả lời sai, downtime, khách không quen dùng bot, dữ liệu không đồng bộ. Với mỗi rủi ro nêu xác suất, mức độ ảnh hưởng, và cách xử lý.

---

## Chương 5 — Đánh giá & Kết luận

**5.1 Tiêu chí đánh giá thành công**
Định nghĩa KPI cụ thể: tỷ lệ booking qua bot, thời gian phản hồi trung bình, tỷ lệ câu hỏi FAQ được giải quyết tự động, no-show rate trước/sau...

**5.2 Hạn chế của đề tài**
Trình bày thẳng thắn: chưa có dữ liệu thực nghiệm, chưa tích hợp thanh toán online, phụ thuộc vào chất lượng knowledge base do owner nhập.

**5.3 Hướng phát triển tiếp theo**
Gợi ý 3–4 hướng mở rộng khả thi: tích hợp thanh toán, voice bot, dự báo doanh thu, chương trình referral.

**5.4 Kết luận**
Tóm tắt đóng góp của đề tài: không chỉ là ứng dụng AI vào một ngành cụ thể, mà là một blueprint có thể replicate cho các chuỗi dịch vụ tương tự (nail, spa, gym...).

---

## Lưu ý về hình ảnh & minh họa nên có

Một số diagram nên chuẩn bị để báo cáo thuyết phục hơn:

1. **Sơ đồ kiến trúc hệ thống** (chương 3) — quan trọng nhất
2. **Flowchart luồng đặt lịch** (UC-B1) — minh họa agent hoạt động cụ thể ra sao
3. **Flowchart HITL** — approve flow cho action nhạy cảm
4. **ERD** — thiết kế dữ liệu
5. **Timeline roadmap** — kế hoạch triển khai theo giai đoạn
6. **Bảng so sánh** giải pháp hiện tại vs. hệ thống đề xuất
