Viewed session.go:57-84

Dưới đây là bản tóm tắt chi tiết file `session.go`, giải thích cách nó hiện thực hóa (implement) bộ tiêu chuẩn quản lý phiên làm việc của **Google ADK-go** bằng **Redis**.

---

### 1. Tổng quan (High-Level)
File này đóng vai trò là một **Redis Storage Provider** cho ADK. Mục tiêu là chuyển đổi các yêu cầu quản lý trạng thái của Agent (như: "Nhớ tên khách hàng", "Lưu lại câu trả lời vừa rồi") thành các thao tác đọc/ghi trên database Redis.

Để làm được điều này, file implement **4 Interface cốt lõi** của ADK:
1.  **`session.Service`**: Bộ não điều khiển (Tạo, Xóa, Lấy Session).
2.  **`session.Session`**: Đại diện cho một phiên làm việc đang tồn tại.
3.  **`session.State`**: Kho chứa dữ liệu (biến số, thông tin ngữ cảnh).
4.  **`session.Events`**: Nhật ký các sự kiện/hội thoại đã xảy ra.

---

### 2. Kiến trúc lưu trữ (Data Organization)
File sử dụng một chiến lược đặt tên Key (Namespace) rất chặt chẽ để quản lý dữ liệu trong Redis:
*   `session:{app}:{user}:{sess}`: Lưu thông tin cơ bản và metadata của phiên.
*   `sessions:{app}:{user}`: Một Redis `SET` chứa danh sách ID của tất cả session thuộc về một user.
*   `events:{app}:{user}:{sess}`: Một Redis `LIST` lưu lịch sử hội thoại (Events).
*   `appstate:{app}`: Lưu dữ liệu dùng chung cho toàn bộ ứng dụng.
*   `userstate:{app}:{user}`: Lưu dữ liệu riêng của từng khách hàng.

---

### 3. Đi sâu vào chi tiết hiện thực (Detailed Implementation)

#### A. Quản lý Trạng thái 3 cấp độ (The 3-Layer State)
Đây là đặc điểm quan trọng nhất của ADK. Khi Agent yêu cầu dữ liệu, file này thực hiện logic:
1.  **Load**: Đọc đồng thời từ 3 key (App, User, Session).
2.  **Merge**: Gộp chúng lại thành một bản đồ dữ liệu duy nhất để Agent sử dụng.
3.  **Persistence**: Khi Agent cập nhật dữ liệu, file phân loại xem dữ liệu đó thuộc cấp độ nào để ghi lại đúng key tương ứng trong Redis.

#### B. Sử dụng Redis Hash cho tính nguyên tử (Atomic Updates)
Thay vì lưu toàn bộ state thành một cục JSON lớn (dễ bị ghi đè dữ liệu nếu có 2 tiến trình cùng chạy), file sử dụng Redis `HASH` (`HSet`/`HGet`). Điều này cho phép cập nhật từng trường dữ liệu nhỏ một cách độc lập và an toàn.

#### C. Cơ chế TTL (Time To Live)
*   Mỗi phiên làm việc đều có một thời hạn (`ttl`).
*   Khi có bất kỳ hoạt động nào (tạo mới, cập nhật), file sẽ gọi `EXPIRE` để gia hạn thời gian sống cho session đó trong Redis.
*   Riêng `AppState` và `UserState` thường được cấu hình để tồn tại lâu hơn hoặc vĩnh viễn (sử dụng `PERSIST`).

#### D. Quản lý sự kiện (Event Streaming)
Các sự kiện hội thoại được lưu dưới dạng một danh sách (List). File hỗ trợ các tính năng như:
*   `NumRecentEvents`: Chỉ lấy N tin nhắn gần nhất để tối ưu bộ nhớ cho Agent.
*   `After`: Chỉ lấy các tin nhắn sau một mốc thời gian cụ thể.

---

### 4. Kết luận: Tại sao nó "đúng chuẩn" ADK?
File này đạt chuẩn ADK vì nó không chỉ đơn thuần là lưu dữ liệu vào Redis, mà nó tái hiện hoàn hảo **hệ thống phân cấp bộ nhớ** của ADK. Nó đảm bảo rằng dù Agent của bạn chạy trên server nào, nó cũng sẽ thấy cùng một "ký ức" giống hệt nhau thông qua Redis backend này.

> [!NOTE]
> Các dòng code `var _ session.Service = (*RedisSessionService)(nil)` ở cuối file chính là cái "chốt chặn" cuối cùng để đảm bảo nếu Google ADK-go thay đổi bất kỳ chuẩn nào, file này sẽ báo lỗi ngay để bạn cập nhật kịp thời.