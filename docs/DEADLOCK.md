### 1. Go chờ bao lâu để báo Deadlock?

Câu trả lời sẽ khiến bạn bất ngờ: **Go không hề có đồng hồ bấm giờ (Timer) để báo Deadlock.**

Nó không chờ 5 giây, 10 giây hay 1 phút. Nó báo Deadlock dựa trên **trạng thái (Status)** của các Goroutine, chứ không dựa trên **thời gian**.

* Nếu bạn `time.Sleep(100 * time.Hour)` (ngủ 100 giờ), Go vẫn vui vẻ đợi và **không bao giờ** báo Deadlock. Tại sao? Vì nó biết sau 100 giờ nữa, Goroutine đó sẽ tỉnh dậy và có thể làm gì đó (như gọi `Done()`).
* Nhưng nếu tất cả Goroutine đều đang đợi Channel hoặc WaitGroup mà **không còn ai đang chạy (Running)** và **không có bộ hẹn giờ (Timer/Sleep)** nào đang kích hoạt, Go sẽ báo Deadlock **ngay lập tức** (chỉ mất vài mili giây để nó nhận ra).

---

### 2. "Deadlock thực sự" là gì trong Go?

Trong Go, Deadlock xảy ra khi hệ thống rơi vào trạng thái **"Bế tắc toàn cục"**.

Hãy tưởng tượng một nhóm người trong một căn phòng:
1.  Ông A bảo: "Tôi đợi ông B đưa chìa khóa".
2.  Ông B bảo: "Tôi đợi ông C đưa chìa khóa".
3.  Ông C bảo: "Tôi đợi ông A đưa chìa khóa".
4.  **Quan trọng:** Không còn ông D nào đang đi tìm chìa khóa cả. Cả 3 ông A, B, C đều đang ngồi im (asleep).

Lúc này, Go Runtime (giống như một vị thần đứng ngoài) nhìn vào căn phòng và thấy: "À, cả lũ đều đang ngủ và không còn ai đang làm việc để tạo ra chìa khóa nữa. Bọn này sẽ ngủ mãi mãi". Thế là nó báo **Fatal Error: Deadlock**.

---

### 3. Go phát hiện Deadlock như thế nào? (Cơ chế bên dưới)

Go có một bộ lập lịch (**Scheduler**) cực kỳ thông minh. Nó quản lý danh sách tất cả các Goroutine và trạng thái của chúng:
* `Running`: Đang chạy.
* `Runnable`: Sẵn sàng chạy nhưng đang đợi CPU.
* `Waiting/Asleep`: Đang đợi (đợi Channel, đợi Mutex, đợi WaitGroup, hoặc đang Sleep).

**Cơ chế kiểm tra:**
Mỗi khi một Goroutine chuyển sang trạng thái `Waiting`, Go Scheduler sẽ kiểm tra:
> "Tổng số Goroutine đang ở trạng thái `Running` hoặc `Runnable` có bằng **0** không?"

* Nếu **bằng 0**: Nghĩa là tất cả mọi người đều đang "há miệng chờ sung", không còn ai đang làm việc để giải cứu ai cả.
* **NGOẠI TRỪ:** Nếu có ít nhất một Goroutine đang `Sleep` (hẹn giờ) hoặc đang đợi I/O (như đợi dữ liệu từ mạng), Go sẽ **không** báo Deadlock. Vì nó hy vọng khi hết giờ Sleep hoặc có dữ liệu mạng về, Goroutine đó sẽ tỉnh lại và giải cứu những người khác.

---

### 4. Phân tích trường hợp `wg.Add(4)` nhưng chỉ có 3 Worker của bạn:

1.  **Giai đoạn 1 (1 giây đầu):**
    * 3 Worker đang `time.Sleep`. Trạng thái là `Waiting` nhưng có bộ hẹn giờ.
    * Main đang `wg.Wait()`. Trạng thái là `Waiting`.
    * Go thấy vẫn còn 3 ông sắp tỉnh dậy -> **Không báo Deadlock.**
2.  **Giai đoạn 2 (Sau 1 giây):**
    * 3 Worker tỉnh dậy, in chữ "XONG", gọi `wg.Done()`.
    * Sau khi gọi `Done()`, 3 Worker này kết thúc (`return`). Chúng biến mất khỏi hệ thống.
    * Lúc này, bộ đếm WaitGroup vẫn còn **1** (vì bạn `Add(4)` mà mới `Done` 3 lần).
    * Hàm `main` vẫn đang `wg.Wait()`.
3.  **Giai đoạn 3 (Khoảnh khắc định mệnh):**
    * Lúc này hệ thống chỉ còn đúng **1 Goroutine duy nhất là Main**.
    * Main đang ở trạng thái `Waiting` (đợi ông thứ 4).
    * Scheduler kiểm tra: "Còn ai đang chạy không? - Không. Còn ai đang Sleep không? - Không. Còn ai đang đợi mạng không? - Không."
    * **Kết luận:** Main sẽ đợi mãi mãi. **BÁO DEADLOCK NGAY LẬP TỨC.**

---

### 💡 Bài học cốt lõi để bạn "suy ra mọi biến thể":

* **Đừng bao giờ để số lượng `Add()` lớn hơn số lượng `Done()` thực tế.**
* Go chỉ báo Deadlock nếu **toàn bộ chương trình** bị kẹt. Nếu bạn có 100 Goroutine, 99 ông bị kẹt nhưng vẫn còn 1 ông đang chạy vòng lặp `for { fmt.Println("Tôi vẫn sống") }`, Go sẽ **không bao giờ** báo Deadlock.
* **Kẹt (Leak)** đáng sợ hơn **Deadlock**. Vì Deadlock làm sập chương trình ngay (bạn biết để sửa), còn Leak (như ví dụ `generator` trước đó) khiến chương trình vẫn chạy nhưng tốn RAM âm thầm, Go không báo lỗi gì cả.

Bạn thấy cơ chế "vị thần đứng ngoài" (Scheduler) này của Go có công bằng không? Nó thà giết nhầm (báo lỗi) còn hơn để chương trình của bạn chạy lãng phí tài nguyên mãi mãi!