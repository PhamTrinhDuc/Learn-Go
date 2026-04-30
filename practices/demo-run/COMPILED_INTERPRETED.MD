## Compiled vs Interpreted — Cách hoạt động

### Interpreted (Python)

Khi bạn chạy `python main.py`, Python không hiểu code của bạn ngay. Nó cần một "phiên dịch viên" đọc từng dòng và dịch sang ngôn ngữ máy **tại thời điểm chạy**:

```
main.py → Python đọc dòng 1 → dịch → CPU chạy
          Python đọc dòng 2 → dịch → CPU chạy
          Python đọc dòng 3 → dịch → CPU chạy
          ...
```

Giống như bạn nói chuyện với người nước ngoài qua phiên dịch viên — **phiên dịch viên phải có mặt suốt buổi nói chuyện**. Không có phiên dịch viên (Python interpreter) thì không nói chuyện được.

---

### Compiled (Go)

Khi bạn chạy `go build`, compiler dịch **toàn bộ code một lần** trước khi chạy, tạo ra file binary là machine code luôn:

```
main.go → compiler dịch toàn bộ → mcp-server.exe (machine code)
                                        ↓
                                   CPU chạy thẳng
                                   (không cần Go nữa)
```

Giống như bạn thuê người dịch **toàn bộ cuốn sách** sang tiếng Việt một lần. Sau đó ai cũng đọc được bản tiếng Việt đó — **không cần phiên dịch viên nữa**.

---

### Tại sao dẫn đến sự khác nhau về Docker image

```
Python image phải mang theo:          Go image chỉ cần:
├── Python interpreter (~30MB)        └── mcp-server.exe (5MB)
├── pip packages (~100MB+)
├── source code (.py files)
└── runtime libs
= 200-400MB                           = 15-20MB
```

Redis viết bằng C cũng cùng lý do — compile xong thành `redis-server` binary ~3MB, Alpine ~5MB, cộng thêm một vài thứ cần thiết → ~60MB.

---

### Một câu để nhớ

> **Interpreted** = phiên dịch viên đi cùng mãi mãi.
> **Compiled** = dịch xong rồi, không cần phiên dịch viên nữa.