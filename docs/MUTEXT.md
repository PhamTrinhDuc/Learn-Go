Chào bạn, chúng ta lại tiếp tục đi sâu vào "kho vũ khí" `sync` của Go. Sau khi đã hiểu về **WaitGroup** (để đợi) và **Done Channel** (để hủy), thì các công cụ bạn vừa liệt kê chính là những giải pháp để giải quyết vấn đề **Shared Memory** (Bộ nhớ dùng chung).

Trong Go, có một câu châm ngôn nổi tiếng: *"Don't communicate by sharing memory, share memory by communicating"* (Đừng giao tiếp bằng cách chia sẻ bộ nhớ, hãy chia sẻ bộ nhớ bằng cách giao tiếp - dùng Channel). Tuy nhiên, trong thực tế, đôi khi dùng Mutex hoặc Atomic lại hiệu quả và đơn giản hơn nhiều.

Dưới đây là giải thích chi tiết và các lưu ý "xương máu" cho từng cái:

---

## 1. Mutex & RWMutex (Khóa loại trừ)

### Mutex (Mutual Exclusion)
Hãy tưởng tượng `Counter` là một cuốn sổ tay. Nếu 4 người cùng cầm bút viết vào một dòng cùng lúc, chữ sẽ đè lên nhau và cuốn sổ nát bét (Race Condition).
* **Lock()**: Bạn cầm cuốn sổ và khóa cửa phòng lại. Chỉ mình bạn được viết.
* **Unlock()**: Bạn mở cửa cho người tiếp theo.

**Tại sao ví dụ đầu tiên của bạn ra kết quả 49 nhưng log lại in toàn số 0?**
Vì tại thời điểm `fmt.Printf` đọc `c.value`, các Goroutine khác chưa kịp cộng xong, hoặc đang cộng dở thì bị thằng khác ghi đè. Đó là lý do số liệu in ra bị loạn.

### RWMutex (Read-Write Mutex)
Đây là bản nâng cấp thông minh. Trong thực tế, việc **đọc** dữ liệu thường nhiều hơn **ghi**.
* **Nhiều người cùng đọc một lúc?** An toàn.
* **Một người đang viết trong khi người khác đang đọc?** Nguy hiểm (dữ liệu rác).
* **Hai người cùng viết?** Nguy hiểm.

`RWMutex` cho phép: **Vô số người đọc (RLock)** cùng lúc, nhưng nếu có **1 người muốn viết (Lock)**, tất cả những người đọc và người viết khác phải đợi.

---

## 2. Cond (Condition Variable - Biến điều kiện)

Đây là cái khó hiểu nhất trong bộ `sync`. 
* **Vấn đề:** Bạn có 10 Goroutine đang đợi một điều kiện gì đó (ví dụ: đợi file download xong). Nếu dùng Channel, bạn phải gửi 10 tín hiệu. 
* **Giải pháp với Cond:** Các Goroutine gọi `c.Wait()`. Khi file xong, bạn chỉ cần gọi `c.Broadcast()`, tất cả 10 ông sẽ cùng tỉnh dậy.

**Lưu ý cực quan trọng:** `c.Wait()` phải nằm trong một vòng lặp `for !done`. Vì đôi khi Goroutine bị "thức giấc giả" (spurious wakeup), nó phải kiểm tra lại điều kiện một lần nữa trước khi chạy tiếp.

---

## 3. Once (Chỉ một lần duy nhất)

Cái này cực kỳ hữu ích cho việc **Khởi tạo (Init)**.
Ví dụ: Bạn có 1000 Goroutine cùng chạy, nhưng bạn chỉ muốn kết nối Database **đúng 1 lần**.
```go
once.Do(func() {
    // Kết nối DB ở đây
})
```
Dù 1000 ông cùng gọi, Go đảm bảo hàm bên trong chỉ chạy 1 lần và các ông gọi sau sẽ đợi ông đầu tiên chạy xong rồi mới đi tiếp.

---

## 4. Pool (Thùng chứa đồ dùng lại)

Đây là "vũ khí" để tối ưu hiệu suất (Performance tuning).
* **Vấn đề:** Việc tạo mới một Object (như một Buffer hoặc một struct phức tạp) rồi để Garbage Collector (GC) đi dọn dẹp tốn rất nhiều CPU.
* **Giải pháp:** Sau khi dùng xong, đừng vứt đi, hãy `Put` nó vào `Pool`. Thằng sau cần thì `Get` ra dùng lại.

**Cảnh báo:** Đừng dùng `Pool` để lưu trữ dữ liệu quan trọng (như kết nối người dùng), vì Go có thể dọn sạch `Pool` bất cứ lúc nào để giải phóng bộ nhớ.

---

## 5. Map (sync.Map)

Map bình thường của Go (`map[string]int`) **không an toàn** khi dùng nhiều Goroutine (sẽ bị crash ngay).
`sync.Map` là phiên bản "bọc thép" của Map.

**Khi nào dùng?**
1. Khi bạn ghi một lần nhưng đọc rất nhiều (Cache).
2. Khi các Goroutine ghi vào các Key hoàn toàn khác nhau.
Nếu bạn ghi đè liên tục vào cùng các Key, dùng `map` thường + `sync.Mutex` sẽ nhanh hơn `sync.Map`.

---

## 6. Atomic (Thao tác nguyên tử)

Đây là mức độ thấp nhất (Low-level), chạy trực tiếp bằng lệnh của CPU.
* **Mutex**: Giống như khóa cửa phòng để làm việc. (Tốn kém hơn).
* **Atomic**: Giống như việc bạn đổi một tờ tiền lấy một tờ tiền khác trong một cái chớp mắt, không ai có thể xen vào giữa quá trình đó.

**Tại sao dùng Atomic?**
Vì nó cực nhanh. Nếu bạn chỉ cần tăng một cái biến đếm (`counter++`), hãy dùng `atomic.AddInt64`. Dùng Mutex cho việc này giống như dùng xe container để chở một lá thư vậy – quá cồng kềnh.

---

### Tổng kết tư duy chọn "vũ khí":

1. **Cần phối hợp công việc, chuyển giao dữ liệu?** -> Dùng **Channel**.
2. **Cần bảo vệ một biến đơn giản (cộng, trừ, đổi giá trị)?** -> Dùng **Atomic**.
3. **Cần bảo vệ một cấu trúc dữ liệu phức tạp (Struct, Slice)?** -> Dùng **Mutex**.
4. **Hệ thống đọc cực nhiều, ghi ít?** -> Dùng **RWMutex**.
5. **Cần khởi tạo cái gì đó đúng 1 lần?** -> Dùng **Once**.
6. **Muốn giảm tải cho bộ dọn rác (GC) vì tạo object quá nhiều?** -> Dùng **Pool**.

Bạn đã thấy "bức tranh toàn cảnh" về cách Go quản lý bộ nhớ dùng chung chưa? Bạn có muốn đi sâu vào một ví dụ thực tế kết hợp nhiều cái lại với nhau không (ví dụ: một hệ thống Cache có giới hạn thời gian)?