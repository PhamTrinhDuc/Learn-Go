---
name: hairchain-agent-system
description: >
  Skill cho dự án thiết kế hệ thống Multi-Agent AI quản lý chuỗi cắt tóc. 
  Kích hoạt khi người dùng hỏi về nghiệp vụ, kiến trúc hệ thống, thiết kế agent, 
  luồng HITL, schema dữ liệu, phân quyền, hoặc bất kỳ khía cạnh kỹ thuật/nghiệp vụ 
  nào liên quan đến dự án này. Cũng kích hoạt khi cần viết thêm chương báo cáo, 
  mở rộng use case, thiết kế API, vẽ diagram, hay đề xuất cải tiến kiến trúc. 
  Dùng skill này ngay cả khi câu hỏi ngắn gọn như "thêm agent mới", "review design", 
  "viết chương 4", hay "luồng thanh toán trông như thế nào" — bất cứ khi nào 
  liên quan đến hệ thống multi-agent chuỗi cắt tóc này.
---
 
# Hairchain — Multi-Agent AI System
 
Dự án thiết kế hệ thống **Multi-Agent AI** cho chuỗi cắt tóc vận hành bởi **1 owner, nhiều chi nhánh**. Mục tiêu: agent xử lý 80% vận hành hàng ngày, owner chỉ can thiệp vào quyết định quan trọng.
 
---
 
## Bối cảnh & Ràng buộc thiết kế
 
- **Doanh nghiệp**: 1 owner, nhiều chi nhánh, có stylist tại mỗi chi nhánh
- **Triết lý**: Agent lo vận hành, Human quyết định rủi ro cao (HITL)
- **Kênh**: Khách dùng Zalo OA, Owner dùng Dashboard + Slack, Manager dùng Mobile App, Stylist dùng app đơn giản xem lịch
---
 
## Kiến trúc hệ thống
 
**Pattern**: Orchestrator → Sub-agent chuyên trách → Storage/Notify
 
```
[Khách/Zalo] [Owner/Dashboard] [Manager/App] [Stylist/App]
                        ↓
                  ORCHESTRATOR
          (intent → route → context → HITL)
                        ↓
   ┌──────┬──────┬──────────┬─────────┬───────────┐
[Booking][FAQ][Inventory][Loyalty][Analytics]
                        ↓
         [PostgreSQL] [Vector DB] [Zalo/Slack API]
```
 
**Nguyên tắc cốt lõi:**
- Single Responsibility — mỗi agent 1 domain
- Stateless Agent — context giữ ở Orchestrator
- Fail-Safe — không chắc → escalate, không tự quyết
- Audit Log — mọi ghi/xóa đều có log + actor + timestamp
---
 
## 5 Sub-Agents — Tóm tắt nhanh
 
| Agent | Đọc | Ghi | Notify |
|---|---|---|---|
| **Booking** | Lịch stylist, slot trống | Tạo/sửa/hủy booking | Zalo (khách), Slack (manager nếu hủy gấp) |
| **FAQ** | Vector Store, bảng giá DB | Không | Slack (câu hỏi không trả lời được) |
| **Inventory** | Tồn kho | Trừ tồn kho bán lẻ | Slack (manager khi thấp) |
| **Loyalty** | Điểm, lịch sử khách | Cộng điểm, tạo voucher | Zalo (nhắc lịch, sinh nhật, điểm) |
| **Analytics** | Toàn bộ dữ liệu | Không | Slack (owner khi bất thường) |
 
---
 
## Booking Agent — Chi tiết nghiệp vụ
 
**Use cases:** UC-B1 (đặt lịch mới), UC-B2 (đổi/hủy), UC-B3 (reminder 24h & 1h tự động), UC-B4 (no-show sau 15 phút)
 
**Quy tắc:**
- Slot = thời lượng dịch vụ + 5 phút buffer
- Không đặt lịch quá 30 ngày tới
- Hủy trong 2h → Slack notify manager
- Hủy > 3 lịch cùng lúc → **bắt buộc HITL owner approve**
**Ngoài scope:** đặt lịch nhóm, thanh toán khi đặt, quản lý ca stylist
 
---
 
## FAQ Agent — Chi tiết nghiệp vụ
 
**Use cases:** UC-F1 (giá/dịch vụ), UC-F2 (địa chỉ/giờ), UC-F3 (upsell tự nhiên), UC-F4 (câu hỏi không có trong KB → Slack manager)
 
**Quy tắc:**
- Giá lấy từ DB theo chi nhánh (không hardcode trong KB)
- Semantic search score < 0.75 → thừa nhận không biết, không bịa
- Upsell tối đa 1 lần / cuộc hội thoại
- Khiếu nại → escalate ngay, không cố xử lý
**Ngoài scope:** tư vấn kỹ thuật tóc chuyên sâu, xử lý khiếu nại, đa ngôn ngữ
 
---
 
## Inventory Agent — Chi tiết nghiệp vụ
 
**Mô hình tồn kho:**
```
Tồn kho thực = Tổng nhập − Bán lẻ (agent trừ) − Nội bộ (manager trừ thủ công)
```
 
**Use cases:** UC-I1 (manager nhập hàng), UC-I2 (khách mua → HITL xác nhận), UC-I3 (cảnh báo ngưỡng thấp → Slack), UC-I4 (tự ngừng bán khi ≤ ngưỡng tối thiểu)
 
**Phân quyền:**
- Nhập hàng / điều chỉnh nội bộ: Owner + Manager
- Trừ tồn kho bán lẻ: chỉ Agent
- Đặt ngưỡng, đổi giá, bật/tắt bán lẻ: chỉ Owner
**Ngoài scope:** tự đặt hàng nhà cung cấp, real-time sync giữa chi nhánh, quản lý hạn sử dụng
 
---
 
## Loyalty Agent — Chi tiết nghiệp vụ
 
**Phân loại khách:** Mới → Thường xuyên → VIP → Ngủ đông (dựa trên tần suất & doanh thu tích lũy, recalculate hàng tuần)
 
**Use cases:**
- UC-L1: Tích điểm sau giao dịch hoàn thành (không tích khi đặt lịch)
- UC-L2: Khách hỏi điểm / đổi điểm lấy voucher (HITL xác nhận trước khi trừ)
- UC-L3: Nhắc quay lại sau 4 tuần, lần 2 sau 6 tuần kèm ưu đãi nhỏ
- UC-L4: Gửi voucher sinh nhật lúc 9:00 sáng ngày sinh nhật
- UC-L5: Recalculate tier hàng tuần, chỉ notify khi thăng hạng
**Quy tắc:** Tỷ lệ điểm do owner cấu hình. Tặng điểm ngoài chương trình → HITL. Khách có thể opt-out bất cứ lúc nào.
 
---
 
## Analytics Agent — Chi tiết nghiệp vụ
 
**Use cases:**
- UC-A1: Báo cáo doanh thu qua ngôn ngữ tự nhiên ("doanh thu tuần này Q1?")
- UC-A2: Phân tích rush hour theo chi nhánh
- UC-A3: Hiệu suất stylist (chỉ owner xem, không ranking công khai)
- UC-A4: Cảnh báo bất thường: doanh thu < 70% baseline, no-show > 30%, 0 booking cả ngày
**Analytics agent chỉ đọc — không bao giờ ghi dữ liệu.**
 
**Ngoài scope:** dự báo doanh thu, so sánh đối thủ, báo cáo thuế/kế toán
 
---
 
## Human-in-the-Loop (HITL)
 
**Action bắt buộc qua owner approve (Slack):**
 
| Action | Timeout |
|---|---|
| Hủy lịch hàng loạt (> 3 lịch) | 30 phút |
| Hoàn tiền cho khách | 1 giờ |
| Tặng điểm/voucher ngoài chương trình | 1 giờ |
| Blast message toàn bộ khách | 24 giờ |
 
**Action bắt buộc qua khách xác nhận (Zalo):**
- Đặt lịch, mua hàng, đổi điểm → timeout 5–10 phút
**Khi timeout:** action bị hủy tự động, gửi Slack nhắc lại.
 
---
 
## Schema dữ liệu — 7 Entity chính
 
**Customer**: customer_id, phone (UNIQUE), name, birth_date, loyalty_points, customer_tier (new/regular/vip/dormant), no_show_count, opt_out_notify, preferred_branch
 
**Branch**: branch_id, name, address, phone, open_time, close_time, slack_channel_id, is_active
 
**Stylist**: stylist_id, branch_id, name, specialties[], is_active
 
**Booking**: booking_id, customer_id, stylist_id, branch_id, service_id, scheduled_at, duration_minutes, status (pending/confirmed/completed/cancelled/no_show), reminder flags, cancelled_by, cancel_reason
 
**Product**: product_id, name, retail_price, unit, is_retail (bool), retail_min_stock, alert_threshold_retail, alert_threshold_internal
 
**Inventory**: product_id × branch_id, total_imported, retail_sold, internal_used, current_stock (computed), is_available
 
**Order**: order_id, customer_id, product_id, branch_id, quantity, unit_price, total_amount, status, points_earned
 
**Audit tables**: StockLog (mọi thay đổi tồn kho), HitlRequest (pending/approved/rejected/timeout)
 
---
 
## Ma trận phân quyền (tóm tắt)
 
| Domain | Owner | Manager | Agent | Khách |
|---|:---:|:---:|:---:|:---:|
| Xem báo cáo toàn chuỗi | ✓ | ✗ | ✓ | ✗ |
| Xem báo cáo chi nhánh | ✓ | ✓ (mình) | ✓ | ✗ |
| Tạo/sửa/hủy booking | ✓ | ✓ | ✓* | ✗ |
| Hủy hàng loạt | ✓** | ✗ | ✗ | ✗ |
| Nhập hàng / điều chỉnh kho | ✓ | ✓ | ✗ | ✗ |
| Trừ tồn kho bán lẻ | ✗ | ✗ | ✓ | ✗ |
| Thay đổi giá / ngưỡng tồn kho | ✓ | ✗ | ✗ | ✗ |
| Tích điểm tự động | ✗ | ✗ | ✓ | ✗ |
| Tặng điểm ngoài chương trình | ✓** | ✗ | ✗ | ✗ |
| Cập nhật knowledge base | ✓ | ✓ | ✗ | ✗ |
 
*`✓*` = Agent thực hiện sau khi khách xác nhận. `✓**` = Qua HITL Slack.*
 
---
 
## Hướng dẫn khi làm việc với skill này
 
### Khi trả lời câu hỏi nghiệp vụ
- Luôn đối chiếu với scope đã định nghĩa (xem mục "Ngoài scope" của từng agent)
- Nếu yêu cầu mới mâu thuẫn với quy tắc hiện tại → nêu rõ conflict và đề xuất cách giải quyết
- Phân quyền phải nhất quán với ma trận phân quyền ở trên
### Khi thiết kế thêm tính năng / agent mới
Trình bày theo template chuẩn:
1. **Mục tiêu & Phạm vi trách nhiệm**
2. **Use cases** (đặt tên UC-X1, UC-X2...)
3. **Input → Output → Tools sử dụng** (dạng code block)
4. **Quy tắc nghiệp vụ** (dạng bảng)
5. **Ngoài scope**
### Khi viết thêm nội dung báo cáo
- Chương 1: Bối cảnh thị trường & pain points → xem GUIDELINE.MD
- Chương 2: Lý thuyết (Agentic AI, HITL, RAG)
- Chương 3: Phân tích & thiết kế → đã viết, dùng làm reference
- Chương 4: Kế hoạch triển khai (MVP → Beta → Full, mỗi giai đoạn 4 tuần)
- Chương 5: Đánh giá & Kết luận
### Khi vẽ diagram (ASCII)
- Kiến trúc tổng thể: Actor → Kênh → Orchestrator → Agents → Storage
- Flowchart: dùng ký hiệu `→`, `↓`, `┌┐└┘│─`, box brackets `[...]`
- HITL flow: luôn có nhánh APPROVE / REJECT / TIMEOUT
---
 
## Công nghệ stack
 
| Layer | Lựa chọn | Lý do |
|---|---|---|
| LLM | Claude (Anthropic) | Tool use, context dài, tiếng Việt tốt |
| Customer channel | Zalo OA API | Phổ biến nhất tại VN cho dịch vụ |
| Internal notify | Slack API | Owner/Manager đã quen dùng |
| Database | PostgreSQL | Quan hệ rõ ràng, audit log dễ |
| Vector DB | Pinecone hoặc Qdrant | FAQ semantic search |
| Hosting | Cloud VM (AWS/GCP) | Dễ scale từng service |
 