Đặc điểm,done Channel (Tín hiệu Đóng),sync.WaitGroup (Bộ đếm)
Hướng điều khiển,"Từ Trên xuống (Top-down): Main bảo Goroutine ""Dừng lại ngay!"".","Từ Dưới lên (Bottom-up): Goroutine bảo Main ""Đợi tôi xong đã!""."
Mục đích,Hủy bỏ (Cancellation): Dùng khi muốn ngắt một việc đang làm dở.,Đợi kết thúc (Waiting): Dùng khi muốn chắc chắn mọi việc đã làm xong 100%.
Sử dụng khi nào?,"Khi bạn làm Pipeline, Stream dữ liệu, hoặc muốn giới hạn thời gian (Timeout).",Khi bạn chia một mảng lớn thành nhiều phần để tính toán song song.
Dữ liệu,Có thể đóng channel để báo hiệu cho hàng triệu Goroutine cùng lúc.,Chỉ là một con số nhảy lên nhảy xuống.