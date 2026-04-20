# AI-Powered Local Service Manager
### Hệ thống Multi-Agent quản lý doanh nghiệp 1 người
**Đề án môn học — Định hướng Đồ án Tốt nghiệp**
Chuyên ngành: Công nghệ Thông tin | Lĩnh vực: Dịch vụ địa phương — AI Agent System

---

## Mục lục

1. [Tổng quan dự án](#1-tổng-quan-dự-án)
2. [Kiến trúc hệ thống](#2-kiến-trúc-hệ-thống)
3. [Thiết kế chi tiết từng Agent](#3-thiết-kế-chi-tiết-từng-agent)
4. [Tech Stack & GenOps](#4-tech-stack--genops)
5. [Data Model / ERD](#5-data-model--erd)
6. [Use Case & Luồng hoạt động](#6-use-case--luồng-hoạt-động)
7. [Wireframe UI](#7-wireframe-ui)
8. [Kế hoạch triển khai & Rủi ro](#8-kế-hoạch-triển-khai--rủi-ro)
9. [Hướng phát triển (Đồ án TN)](#9-hướng-phát-triển-đồ-án-tốt-nghiệp)

---

## 1. Tổng quan dự án

### 1.1 Bối cảnh & Vấn đề

Chủ tiệm dịch vụ địa phương (cắt tóc, nail, sửa xe, quán ăn nhỏ...) thường vận hành **một mình** — vừa làm dịch vụ, vừa nhận đặt lịch, vừa quản lý kho, vừa theo dõi tài chính. Đây là bài toán *"một tay làm tất"* điển hình.

**Các điểm đau cốt lõi:**
- Bỏ lỡ lịch hẹn khi đang bận tay phục vụ khách
- Không theo dõi được tồn kho — hết vật tư lúc đang làm
- Không có báo cáo doanh thu — mọi thứ ước tính trong đầu
- Khách hỏi giá/chính sách phải dừng tay trả lời hoặc bỏ qua
- Không có thời gian làm marketing — chỉ truyền miệng

### 1.2 Giải pháp

Xây dựng hệ thống **Multi-Agent AI** tích hợp vào website, hoạt động 24/7 như một trợ lý ảo toàn năng:
- Tự động xử lý các tác vụ vận hành thường ngày
- Thông báo cho chủ qua **Slack**
- Chờ approval với các tác vụ quan trọng (risk-based)

### 1.3 Phạm vi hệ thống

| Thành phần | Mô tả | Phase |
|---|---|---|
| Website (UI) | Giao diện khách hàng đặt lịch, hỏi thông tin | MVP |
| Orchestrator Agent | Não trung tâm — route intent, quản lý approval flow | MVP |
| Booking Agent | Nhận/xác nhận/nhắc lịch hẹn tự động | MVP |
| FAQ/RAG Agent | Trả lời giá cả, chính sách bằng RAG | MVP |
| Inventory Agent | Theo dõi tồn kho, cảnh báo hết hàng | MVP |
| Analytics Agent | Báo cáo doanh thu, xu hướng theo lịch | Phase 2 |
| Slack Notification | Kênh giao tiếp chủ tiệm — approval & alerts | MVP |

---

## 2. Kiến trúc hệ thống

### 2.1 Sơ đồ tổng thể

```
┌──────────────────────────────────────────────────────────────────┐
│                        INTERFACE LAYER                           │
│         Website (Customer-facing)    Slack (Owner-facing)        │
└──────────────────────────┬───────────────────┬───────────────────┘
                           │ request/response  │ notify/approval
                           ▼                   ▼
┌──────────────────────────────────────────────────────────────────┐
│                     ORCHESTRATOR AGENT                           │
│     Intent Routing | Approval Flow Manager | Memory Access       │
└───────────┬────────────────────────────────────────┬────────────┘
            │ coordinate                             │ get memory
            ▼                                        ▼
┌───────────────────────────────┐      ┌─────────────────────────────┐
│         SUB-AGENT LAYER       │      │       MEMORY LAYER           │
│  ┌─────────┐  ┌───────────┐  │      │  In-Context: LLM window      │
│  │ Booking │  │  FAQ/RAG  │  │      │  Short-term: PostgreSQL       │
│  │  Agent  │  │   Agent   │  │      │  Long-term: PGVector (RAG)   │
│  └────┬────┘  └─────┬─────┘  │      └─────────────────────────────┘
│  ┌────┴────┐  ┌─────┴─────┐  │
│  │Inventory│  │ Analytics │  │
│  │  Agent  │  │   Agent   │  │
│  └────┬────┘  └─────┬─────┘  │
└───────┼─────────────┼────────┘
        │    call     │
        ▼             ▼
┌──────────────────────────────────────────────────────────────────┐
│                         TOOLS LAYER                              │
│  check_slots | create_booking | semantic_search | check_inventory│
│  update_stock | get_analytics | send_notification | ...          │
└───────────────────────────┬──────────────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────────┐
│                        DATABASE LAYER                            │
│          PostgreSQL (RDBMS)          PGVector (Vector Store)     │
└──────────────────────────────────────────────────────────────────┘
```

### 2.2 Orchestrator Agent — Não trung tâm

Orchestrator nhận mọi input (từ website hoặc Slack), xác định intent, quyết định agent nào xử lý, và quản lý toàn bộ approval flow.

| Trách nhiệm | Chi tiết |
|---|---|
| Intent Routing | Hybrid: Rule-based pre-filter cho intent đơn giản → LLM reasoning cho intent mơ hồ/phức tạp |
| Approval Flow | Phân loại action theo risk level: Auto / Soft-confirm / Hard-approval |
| Memory Management | Quản lý 3 tầng memory: In-context, Short-term (PG), Long-term (PGVector) |
| Agent Coordination | Gọi sub-agent qua tool calling, tổng hợp kết quả, trả response |
| Slack Integration | Gửi notify + nhận approval callback từ chủ tiệm |

**Routing logic (hybrid):**
```
Input nhận được
    │
    ▼
[Rule-based pre-filter]
    │ Keyword/pattern match rõ ràng?
    ├── YES → Route thẳng đến agent (low latency, no LLM)
    │         Ví dụ: "đặt lịch lúc 3h" → Booking Agent
    │
    └── NO  → LLM Orchestrator reasoning
              Ví dụ: "hôm nay tôi có bận không?" → cần query lịch + reason
```

### 2.3 Approval Flow Design

Điểm mấu chốt về UX — chủ tiệm thường bận tay, không thể approve ngay. Hệ thống phân loại 3 mức risk:

| Risk Level | Ví dụ action | Behavior |
|---|---|---|
| 🟢 Auto-run | Trả lời FAQ, gửi nhắc lịch, log thay đổi kho | Thực thi ngay, notify sau (async) |
| 🟡 Soft confirm | Tạo lịch hẹn mới, cảnh báo tồn kho thấp | Notify Slack → tự chạy sau 5 phút nếu không phản hồi |
| 🔴 Hard approval | Hủy lịch, xóa dữ liệu, thay đổi giá | Bắt buộc chủ bấm Approve/Reject trên Slack Block Kit |

### 2.4 Memory Architecture

| Tầng | Backend | Lưu gì | Khi nào dùng |
|---|---|---|---|
| In-Context (Working) | LLM context window | Cuộc hội thoại hiện tại, agent state | Luôn luôn |
| Short-term (Episodic) | PostgreSQL | Lịch sử booking, inventory logs, 30 ngày gần nhất | Query trực tiếp, không qua RAG |
| Long-term (Semantic) | PGVector | Profile khách, knowledge base tiệm, giá/dịch vụ/FAQ | Chỉ trigger khi query liên quan khách cụ thể hoặc KB tiệm |

> **Design note:** RAG không được dùng cho mọi thứ. Booking và Finance query thẳng PostgreSQL — nhanh hơn, deterministic hơn. RAG chỉ dành cho unstructured knowledge.

---

## 3. Thiết kế chi tiết từng Agent

### 3.1 Booking Agent

**Metadata:**

| | |
|---|---|
| Serve | Customer-facing |
| Trigger | On-demand — khi khách đặt/hỏi/hủy lịch |
| Risk mặc định | 🟡 Soft confirm (tạo lịch mới) |
| LLM | GPT-4o mini với tool calling |

**Tool set:**

| Tool | Input | Output |
|---|---|---|
| `check_available_slots` | date, service_type, duration | Danh sách slot trống |
| `create_booking` | customer_id, slot_id, service, notes | Booking ID, confirmation |
| `cancel_booking` | booking_id, reason | Status (🔴 cần Hard approval) |
| `send_reminder` | booking_id, channel | Delivery status |
| `get_customer_history` | customer_id | Lịch sử đặt lịch, preferences |

**Agentic Loop:**
```
Khách: "Đặt lịch cắt tóc sáng mai 9h"
    │
    ▼ Orchestrator route → Booking Agent
    │
    ▼ check_available_slots(tomorrow, 09:00, haircut, 45min)
    │
    ├── Có slot → create_booking() → Slack notify (🟡 soft confirm)
    │                → 5 phút không reject → confirm với khách
    │                → schedule reminder T-1h
    │
    └── Không có slot → Gợi ý 3 slot gần nhất để khách chọn
```

---

### 3.2 FAQ/RAG Agent

**Metadata:**

| | |
|---|---|
| Serve | Customer-facing |
| Trigger | On-demand — khi khách hỏi thông tin |
| Risk mặc định | 🟢 Auto-run (read-only) |

**RAG Pipeline:**

```
Query của khách
    │
    ▼ Embed query (text-embedding-3-small)
    │
    ▼ Similarity search PGVector (top-k=3, threshold=0.75)
    │
    ├── score >= 0.75 → Rerank nếu cần → Generate response với context
    │
    └── score < 0.75 → Fallback: "Tôi chưa có thông tin này,
                        để tôi hỏi chủ tiệm nhé" → Notify Slack
```

**Knowledge Base gồm:**
- Bảng giá từng dịch vụ
- Chính sách đặt lịch / hủy lịch / hoàn tiền
- FAQ thường gặp
- Giờ mở cửa, địa chỉ, liên hệ
- Mô tả chi tiết từng dịch vụ

**Chunking strategy:** Semantic chunking theo từng dịch vụ/chính sách (không dùng sliding window để tránh context bị cắt giữa chừng).

**Tool set:**

| Tool | Input | Output |
|---|---|---|
| `semantic_search` | query, top_k, threshold | Relevant chunks + similarity score |
| `get_service_pricing` | service_name | Giá, mô tả, thời gian thực hiện |
| `get_policy` | policy_type | Nội dung chính sách đầy đủ |

---

### 3.3 Inventory Agent

**Metadata:**

| | |
|---|---|
| Serve | Owner-facing |
| Trigger | Event-driven (sau mỗi booking hoàn thành) + Schedule (8h sáng hàng ngày) |
| Risk mặc định | 🟢 Auto-run (alert) \| 🔴 Hard approval (xóa/sửa dữ liệu lớn) |

**Tool set:**

| Tool | Input | Output |
|---|---|---|
| `get_inventory_status` | item_id? (optional) | Danh sách tồn kho + cảnh báo items sắp hết |
| `update_stock` | item_id, quantity, reason | Updated stock + log entry (🔴) |
| `deduct_from_booking` | booking_id, materials_used | Auto-deduct theo service template (🟢) |
| `alert_low_stock` | threshold config | Slack notification tự động |
| `suggest_reorder` | — | Gợi ý nhập hàng dựa trên usage pattern 7 ngày |

**Luồng event-driven:**
```
Booking hoàn thành
    │
    ▼ Inventory Agent nhận event
    │
    ▼ deduct_from_booking(booking_id) → trừ vật tư theo service template
    │
    ▼ check threshold → nếu bất kỳ item < min_stock
    │
    └── alert_low_stock() → Slack: "Vật tư X còn N đơn vị, cần nhập thêm"
                         → suggest_reorder() đính kèm
```

---

### 3.4 Analytics Agent

**Metadata:**

| | |
|---|---|
| Serve | Owner-facing |
| Trigger | Schedule thuần túy — cron job |
| Risk mặc định | 🟢 Auto-run (read-only + gửi báo cáo) |

**Lịch chạy:**

| Tần suất | Thời gian | Nội dung báo cáo |
|---|---|---|
| Hàng ngày | 22:00 | Tổng lịch hẹn, doanh thu ước tính, no-show rate, khách mới |
| Hàng tuần | CN 20:00 | Khách lâu không quay lại (>30 ngày), giờ cao điểm trong tuần |
| Hàng tháng | Ngày 1 | Dịch vụ bán chạy/chậm, top khách, so sánh tháng trước |

---

## 4. Tech Stack & GenOps

### 4.1 Core Stack

| Layer | Technology | Lý do chọn |
|---|---|---|
| LLM | GPT-4o mini (OpenAI API) | Chi phí thấp, đủ nhanh cho tool calling, JSON mode tốt |
| Orchestration | LangGraph | Native multi-agent, state machine, human-in-the-loop built-in |
| RAG / Vector | pgvector (PostgreSQL extension) | Tích hợp chung Postgres, tránh thêm infra phức tạp |
| Database | PostgreSQL | RDBMS chính cho booking, inventory, customer data |
| Backend | FastAPI (Python) | Async, lightweight, tích hợp tốt với LangGraph |
| Frontend | Next.js / React | SSR tốt cho SEO, dễ build chat UI |
| Notification | Slack Bot API + Block Kit | Approval UI với button, không cần Zalo OAuth2 doanh nghiệp |
| Scheduler | APScheduler (Python) | Embedded cron cho Analytics Agent, không cần Celery |

### 4.2 GenOps Stack

| Mục tiêu | Tool | Scope |
|---|---|---|
| Tracing & Observability | Langfuse | Trace mọi LLM call, tool call, latency, cost/token per agent |
| Structured Logging | Python logging (JSON format) | Ghi log agent action, approval event, error với context |
| Cost Monitoring | Langfuse dashboard | Cost theo ngày, breakdown per agent |
| Deployment | Docker + Docker Compose | Containerize backend, scheduler, PG |
| Hosting | VPS (Railway / DigitalOcean) | Không cần K8s ở MVP phase |
| CI/CD | GitHub Actions | Auto-deploy khi push main branch |
| Secret Management | .env + python-dotenv | API keys, DB credentials |

### 4.3 Lý do chọn LangGraph thay vì LangChain thuần

LangGraph phù hợp hơn vì:
- **State machine native** — approval flow "chờ 5 phút" cần persistent state giữa các bước, LangGraph handle được cleanly
- **Human-in-the-loop** — `interrupt()` node built-in cho hard approval cases
- **Multi-agent coordination** — mỗi sub-agent là 1 node trong graph, Orchestrator là supervisor node
- **Async-first** — chủ tiệm có thể gửi nhiều request đồng thời (khách đặt lịch trong khi agent đang chờ approval)

---

## 5. Data Model / ERD

### 5.1 Các bảng chính (PostgreSQL)

```sql
-- Khách hàng
customers (
    id UUID PK,
    name VARCHAR,
    phone VARCHAR UNIQUE,
    created_at TIMESTAMP,
    last_visit TIMESTAMP,
    notes TEXT                  -- ghi chú thủ công của chủ
)

-- Dịch vụ
services (
    id UUID PK,
    name VARCHAR,
    duration_minutes INT,
    price DECIMAL,
    materials JSONB             -- {"shampoo": 50, "dye": 100} (ml/g)
)

-- Lịch hẹn
bookings (
    id UUID PK,
    customer_id UUID FK → customers,
    service_id UUID FK → services,
    scheduled_at TIMESTAMP,
    status ENUM(pending, confirmed, completed, cancelled, no_show),
    notes TEXT,
    created_at TIMESTAMP,
    approved_by_owner BOOLEAN DEFAULT FALSE
)

-- Tồn kho
inventory_items (
    id UUID PK,
    name VARCHAR,
    unit VARCHAR,               -- ml, g, cái, hộp
    current_stock DECIMAL,
    min_stock DECIMAL,          -- ngưỡng cảnh báo
    updated_at TIMESTAMP
)

-- Log thay đổi kho
inventory_logs (
    id UUID PK,
    item_id UUID FK → inventory_items,
    change_amount DECIMAL,      -- âm = giảm, dương = nhập
    reason VARCHAR,             -- booking_deduct | manual | restock
    booking_id UUID FK → bookings (nullable),
    created_at TIMESTAMP
)

-- Approval queue (cho Orchestrator)
approval_queue (
    id UUID PK,
    action_type VARCHAR,
    payload JSONB,
    risk_level ENUM(soft, hard),
    status ENUM(pending, approved, rejected, auto_executed),
    slack_message_ts VARCHAR,   -- để edit message sau khi xử lý
    created_at TIMESTAMP,
    resolved_at TIMESTAMP
)

-- Conversation history (short-term memory)
conversations (
    id UUID PK,
    session_id VARCHAR,
    role ENUM(user, assistant, tool),
    content TEXT,
    metadata JSONB,
    created_at TIMESTAMP
)
```

### 5.2 PGVector — Knowledge Base

```sql
-- Chunks cho RAG
knowledge_chunks (
    id UUID PK,
    source_type VARCHAR,        -- pricing | policy | faq | info
    source_name VARCHAR,        -- tên dịch vụ / tên chính sách
    content TEXT,
    embedding VECTOR(1536),     -- text-embedding-3-small
    metadata JSONB,
    created_at TIMESTAMP
)

-- Index cho similarity search
CREATE INDEX ON knowledge_chunks
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
```

---

## 6. Use Case & Luồng hoạt động

### 6.1 Actor & Use Case tổng quan

**Actors:**
- **Customer** — khách hàng tương tác qua website
- **Owner** — chủ tiệm tương tác qua Slack
- **System** — các agent tự động theo lịch

**Use cases chính:**

| Actor | Use Case |
|---|---|
| Customer | Đặt lịch hẹn mới |
| Customer | Xem/hủy lịch hẹn |
| Customer | Hỏi giá dịch vụ |
| Customer | Hỏi chính sách, giờ mở cửa |
| Owner | Xem danh sách lịch hẹn hôm nay |
| Owner | Approve/Reject action từ Agent |
| Owner | Cập nhật tồn kho thủ công |
| Owner | Xem báo cáo doanh thu |
| System | Gửi nhắc lịch cho khách (T-1h) |
| System | Trừ vật tư sau booking hoàn thành |
| System | Gửi báo cáo cuối ngày/tuần/tháng |
| System | Cảnh báo tồn kho thấp |

### 6.2 Sequence Diagram — Khách đặt lịch

```
Customer    Website     Orchestrator    Booking Agent    Slack (Owner)
    │           │               │               │               │
    │──đặt lịch►│               │               │               │
    │           │──────────────►│               │               │
    │           │               │──route──────►│               │
    │           │               │               │──check_slots──►│ (DB)
    │           │               │               │◄──available────│
    │           │               │               │──create_booking►│ (DB)
    │           │               │◄──booking_id──│               │
    │           │               │──────────────────────────────►│
    │           │               │               │          notify: "Lịch mới
    │           │               │               │          X đặt lịch 9h mai"
    │           │◄──"Đặt lịch thành công"───────│               │
    │◄──confirm──│               │               │               │
    │           │               │               │         [5 phút...]
    │           │               │               │   không reject → auto confirm
    │           │               │               │──send_reminder─►│ (T-1h)
```

### 6.3 Sequence Diagram — Hard Approval (hủy lịch)

```
Owner       Slack       Orchestrator    Booking Agent    Database
  │           │               │               │               │
  │──"hủy lịch booking #123"─►│               │               │
  │           │               │──route──────►│               │
  │           │               │               │──get_booking──►│
  │           │               │               │◄──booking data─│
  │           │               │◄──action ready│               │
  │           │◄──Block Kit: "Xác nhận hủy lịch #123?"────────│
  │◄──prompt──│               │               │               │
  │──[Approve]►│               │               │               │
  │           │──callback─────►│               │               │
  │           │               │──────────────►│               │
  │           │               │               │──cancel_booking►│
  │           │               │               │◄──success──────│
  │           │◄──"Đã hủy lịch #123"──────────│               │
  │◄──confirm──│               │               │               │
```

---

## 7. Wireframe UI

### 7.1 Website — Trang chủ / Chat Interface

```
┌─────────────────────────────────────────────────────────┐
│  [Logo Tiệm]                        Đặt lịch  Giá cả   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│   ╔═══════════════════════════════════════════╗         │
│   ║  Xin chào! Tôi là trợ lý của tiệm.       ║         │
│   ║  Tôi có thể giúp bạn:                     ║         │
│   ║  • Đặt lịch hẹn                           ║         │
│   ║  • Xem giá dịch vụ                        ║         │
│   ║  • Trả lời câu hỏi về tiệm               ║         │
│   ╚═══════════════════════════════════════════╝         │
│                                                         │
│   [Đặt lịch ngay]  [Xem bảng giá]  [Giờ mở cửa]       │
│                                                         │
│   ┌─────────────────────────────────────────┐           │
│   │ Nhập tin nhắn...                   [Gửi]│           │
│   └─────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────┘
```

### 7.2 Website — Booking Flow

```
┌─────────────────────────────────────────────────────────┐
│  Chọn dịch vụ                                           │
│  ○ Cắt tóc nam (45 phút) — 80,000đ                     │
│  ○ Gội đầu + cắt (60 phút) — 120,000đ                  │
│  ○ Nhuộm tóc (120 phút) — 250,000đ                     │
├─────────────────────────────────────────────────────────┤
│  Chọn ngày & giờ                                        │
│  [◄ Tháng 4, 2026 ►]                                   │
│   T2   T3   T4   T5   T6   T7   CN                     │
│    6    7   [8]   9   10   11   12                      │
│                                                         │
│  Giờ trống ngày 8/4:                                   │
│  [9:00] [10:00] [14:00] [15:30] [16:30]                │
├─────────────────────────────────────────────────────────┤
│  Thông tin của bạn                                      │
│  Họ tên: [________________]                             │
│  SĐT:    [________________]                             │
│  Ghi chú: [________________]                            │
│                                [Xác nhận đặt lịch]     │
└─────────────────────────────────────────────────────────┘
```

### 7.3 Slack — Approval Message (Block Kit)

```
┌─────────────────────────────────────────────────────────┐
│ 🗓️ Lịch hẹn mới                                        │
│─────────────────────────────────────────────────────────│
│ Khách: Nguyễn Văn A (0912 345 678)                     │
│ Dịch vụ: Cắt tóc nam                                   │
│ Thời gian: Thứ 4, 08/04 lúc 9:00 (45 phút)            │
│ Ghi chú: —                                              │
│─────────────────────────────────────────────────────────│
│ ⏱️ Tự động xác nhận sau 5 phút nếu không phản hồi      │
│─────────────────────────────────────────────────────────│
│  [✅ Xác nhận ngay]      [❌ Từ chối]                   │
└─────────────────────────────────────────────────────────┘
```

### 7.4 Slack — Báo cáo cuối ngày (Analytics)

```
┌─────────────────────────────────────────────────────────┐
│ 📊 Tổng kết ngày 08/04/2026                            │
│─────────────────────────────────────────────────────────│
│ ✅ Lịch hoàn thành:    6 / 7                           │
│ ❌ No-show:            1 (Trần Thị B — 14:00)          │
│ 💰 Doanh thu ước tính: 620,000đ                        │
│ 👤 Khách mới:          2                               │
│─────────────────────────────────────────────────────────│
│ 📦 Tồn kho cần chú ý:                                  │
│   • Dầu gội (còn 200ml) — dưới ngưỡng tối thiểu       │
│─────────────────────────────────────────────────────────│
│ [Xem chi tiết]   [Đặt hàng vật tư]                    │
└─────────────────────────────────────────────────────────┘
```

---

## 8. Kế hoạch triển khai & Rủi ro

### 8.1 MVP Scope (Đề án môn học)

Thiết kế đầy đủ 5 agent, nhưng MVP thực thi tập trung vào:

| Priority | Component | Lý do |
|---|---|---|
| P0 | Orchestrator + Booking Agent | Core business value — lịch hẹn là thu nhập |
| P0 | FAQ/RAG Agent | Giảm tải trả lời cho chủ, chạy 24/7 |
| P0 | Slack notification + Approval flow | Chủ tiệm phải kiểm soát được |
| P1 | Inventory Agent | Quan trọng nhưng không critical như booking |
| P2 | Analytics Agent | Cron job — có thể thêm sau khi core ổn |

### 8.2 Rủi ro & Giải pháp

| Rủi ro | Mức độ | Giải pháp |
|---|---|---|
| Race condition trong approval flow — agent đang chờ approve nhưng có request mới vào | Cao | LangGraph state machine + persistent state trong DB, mỗi approval_queue record độc lập |
| LLM hallucinate sai thông tin giá/chính sách | Trung bình | FAQ/RAG Agent dùng retrieved context + system prompt "chỉ trả lời dựa trên context, không sáng tạo" |
| Slack downtime → chủ không nhận được approval request | Thấp | Soft confirm có timeout 5 phút, hard approval có deadline 24h sau đó escalate qua email |
| Cost LLM vượt dự kiến | Trung bình | Rule-based pre-filter giảm 70-80% LLM calls, dùng GPT-4o mini thay GPT-4o |
| Khách spam booking | Thấp | Rate limit theo phone number, require xác nhận OTP cho lần đầu |

---

## 9. Hướng phát triển (Đồ án Tốt nghiệp)

Đề án môn học này là **thiết kế nền tảng** — đồ án tốt nghiệp sẽ implement và mở rộng:

### 9.1 Tách Notification Agent

Hiện tại Slack notification nằm trong Orchestrator. Khi mở rộng multi-channel (email, SMS, Zalo sau khi có tài khoản doanh nghiệp), cần tách thành agent riêng:

```
Notification Agent
    ├── Slack (approval, alerts)
    ├── Email (báo cáo tháng, khách mới)
    ├── SMS (reminder khách hàng)
    └── Zalo OA (tương lai — khi có doanh nghiệp đăng ký)
```

### 9.2 Multi-tenant

Từ 1 tiệm → SaaS cho nhiều tiệm. Cần thêm:
- Tenant isolation ở DB layer (row-level security)
- Per-tenant knowledge base và config
- Billing theo usage (LLM cost tracking per tenant)

### 9.3 Voice Input

Thêm Whisper STT cho chủ tiệm gửi voice note qua Slack thay vì phải gõ — thực tế hơn nhiều khi đang bận tay.

### 9.4 Proactive Customer Engagement

Analytics Agent mở rộng thêm khả năng:
- Tự động nhắn khách lâu không quay lại (>30 ngày) với ưu đãi
- Personalized recommendation dựa trên lịch sử dịch vụ
- A/B testing nội dung marketing

---

*Tài liệu này là bản thiết kế tổng hợp — cập nhật liên tục trong quá trình phát triển.*
*Last updated: 04/2026*