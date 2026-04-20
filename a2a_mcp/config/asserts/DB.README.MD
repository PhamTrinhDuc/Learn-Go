## 1. Cơ chế Bảo mật Multi-tenant (Row-Level Security - RLS)

Đây là phần quan trọng nhất giúp ngăn chặn việc rò rỉ dữ liệu giữa các khách hàng (tenants).

* **Bật RLS:** Lệnh `ALTER TABLE documents ENABLE ROW LEVEL SECURITY;` bắt buộc mọi truy vấn phải đi qua một bộ lọc bảo mật.
* **Chính sách cách ly (Isolation Policy):**
    ```sql
    USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid)
    ```
    * **Cơ chế:** Thay vì dùng `WHERE tenant_id = ?` thủ công trong code Go, hệ thống sử dụng một biến phiên (session variable) của Postgres.
    * **Hoạt động:** Trước mỗi câu query, ứng dụng Go phải chạy lệnh `SET app.current_tenant_id = 'id-cua-khach-hang'`. Postgres sẽ tự động lọc bỏ tất cả các dòng không thuộc về ID này.
* **User Phân quyền (app_user):**
    * Tạo `app_user` riêng biệt với quyền `NOSUPERUSER` và `NOBYPASSRLS`.
    * Việc dùng một User không có quyền bypass RLS giúp đảm bảo ngay cả khi lập trình viên quên thêm điều kiện lọc trong code, database vẫn bảo vệ dữ liệu tuyệt đối.

---

## 2. Vector Search & AI Optimization (pgvector)

Hệ thống được thiết kế đặc biệt để lưu trữ và tìm kiếm vector (thường dùng cho Chatbot AI, RAG).

* **Extension:** Sử dụng `vector` để hỗ trợ kiểu dữ liệu mảng số thực đặc biệt.
* **Định dạng:** `vector(1536)` được tối ưu cho model **OpenAI ada-002**.
* **Index IVFFlat:**
    ```sql
    CREATE INDEX ... USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
    ```
    * Đây là thuật toán phân cụm. Nó chia các vector vào 100 "danh sách" (lists). Khi tìm kiếm, nó chỉ quét các danh sách gần nhất thay vì toàn bộ bảng, giúp tốc độ search nhanh gấp hàng trăm lần.

---

## 3. Chiến lược Indexing Đa lớp (Performance)

File cấu hình này sử dụng 4 loại Index khác nhau cho 4 mục đích riêng biệt:

1.  **B-Tree Index (`idx_documents_tenant_id`):** Tối ưu cho việc lọc nhanh dữ liệu theo khách hàng.
2.  **GIN Index cho Metadata (`idx_documents_metadata`):** Tối ưu cho việc tìm kiếm bên trong cột JSONB. Cho phép bạn tìm các tags, category trong metadata cực nhanh.
3.  **GIN Index cho Full-Text Search (`idx_documents_fulltext`):** * Sử dụng `to_tsvector('english', title || ' ' || content)`.
    * Hỗ trợ tìm kiếm từ khóa (BM25) dựa trên gốc từ (Stemming) và loại bỏ từ dừng (Stop words).
4.  **IVFFlat Index:** Như đã nói ở trên, dùng cho tìm kiếm vector tương đồng.

---

## 4. Tự động hóa Dữ liệu (Trigger & Functions)

Dùng để đảm bảo tính toàn vẹn và theo dõi dấu vết dữ liệu mà không cần can thiệp từ tầng Application.

* **Tự động cập nhật thời gian:** Hàm `update_updated_at_column()` và các `TRIGGER` tương ứng đảm bảo cột `updated_at` luôn phản ánh đúng thời điểm cuối cùng dòng dữ liệu bị thay đổi.
* **Cascade Delete:** Ràng buộc `ON DELETE CASCADE` trên `tenant_id` đảm bảo khi một khách hàng bị xóa, toàn bộ documents và logs liên quan sẽ bị xóa sạch, tránh dữ liệu rác (orphan data).

---

## 5. Quản lý Tài nguyên & Cost Control (Usage Logs)

Thiết kế bảng `usage_logs` cho thấy hệ thống có khả năng quản lý chi phí AI:
* Theo dõi `prompt_tokens`, `completion_tokens` và `total_tokens`.
* Tính toán `cost_usd` trực tiếp ở mức database.
* Index phức hợp `(tenant_id, user_id, created_at DESC)` cho phép truy xuất báo cáo sử dụng của một khách hàng theo thời gian thực rất nhanh.

---

### Tóm tắt luồng hoạt động chuẩn Security:
1.  App Go kết nối vào DB bằng **app_user**.
2.  Mở Transaction.
3.  Chạy lệnh: `SET LOCAL app.current_tenant_id = '...'`.
4.  Thực hiện câu lệnh SQL (Postgres sẽ tự áp dụng RLS).
5.  Commit Transaction.
