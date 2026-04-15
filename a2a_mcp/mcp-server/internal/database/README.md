### I. Khởi tạo connection Pool với PostgreSQL trong Golang 
#### 1. Tạo struct DB 
- Struct chứa tham số pool có kiểu con trỏ bởi sau này khi dùng lại pool thì Go sẽ dùng chính xác pool thông qua địa chỉ thay vì tạo ra bản sao mới nếu không có con trỏ
```bash
type DB struct {
	pool *pgxpool.Pool
}
```
#### 2. Hàm khởi tạo Constructor 
##### Input: 
- Nhận context để hủy toàn bộ luồng request nếu trước đó có lỗi 
- Nhận config bao gồm các thông số cần thiết để kết nối tới PG 
```bash
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32 // tối đa X connect 1 lúc
	MinConns int32 // tối thiếu X connect 1 lúc
}
```
##### Body: 
1. Nối các config thành chuỗi
2. Parse chuỗi này để kiểm tra xem có lỗi trước khi kết nối không => return `poolConfig`
3. Cấu hình các thông số cho `poolConfig`
- Ghi đè lại 2 tham số MaxConns và MinConns 
- Cấu hình thêm 3 tham số thời gian khác
```bash
poolConfig.MaxConns = cfg.MaxConns             // override connString
poolConfig.MinConns = cfg.MinConns             // override connString
poolConfig.MaxConnLifeTime = time.Hour         // thời gian sống tối đa của 1 kết nối
poolConfig.MaxConnIdleTime = 30 * time.Minute  // thời gian tối đa của 1 kết nối không hoạt động
poolConfig.HealthCheckPeriod = 1 * time.Minute // thời gian kiểm tra kết nối
```
- `poolConfig.AfterConnect`: config này nhằm mục đích thực hiện các công việc ngay sau khi kết nối và trước khi bắt đầu thực hiện CURD. Cụ thể ở đây là đăng ký cho PG biết kiểu `vector` là gì nhằm mục đích cấu hình cho vector search. Tuy nhiên hiện nay PG đã biết kiểu vector này nên không cần implement cụ thể, trong code chỉ mock để giả lập. Ngoài hàm này ta có thể thực hiện các hàm khác phục vụ cho công việc khác
```bash
// Implement cụ thể nếu PG chưa biết kiểu vector là gì
poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
    // 1. Kết nối vừa mở xong
    // 2. Ta ra lệnh cho kết nối này đăng ký kiểu dữ liệu Vector
    // 3. Nếu đăng ký lỗi, kết nối này coi như hỏng (return err)
    return pgvector.RegisterTypes(ctx, conn) // Ví dụ thực tế sẽ như thế này
}
```
4. Tạo connection pool 
5. Ping thử để test connection 
##### Return: 
- Trả về địa chỉ của Struct DB, sau này khi dùng db thì go biết đang dùng chính xác instance vừa tạo chứ không tạo ra bản sao của struct này

### II. Bọc Transaction cho Row Level Security
1. Row Level Security (method: setTententContext)
- Khi query dữ liệu, thay vì phải `where user_id = 'ABC'` (có thể quên ở 1 lần gọi nào đó) ta thực hiện: `SET LOCAL app.current_tenant_id`
- Khi 1 connection trong pool được sử dụng (1 user request), db sẽ tự biết rằng connection này chỉ được dùng dữ liệu của `tenant_id` được SET trước đó.
- `app.current_tenant_id`: Postgres cho phép tạo ra các biến tạm thời của riêng mình theo cú pháp:
`SET <tên_nhóm>.<tên_biến> = 'giá trị'`
    + `app`: Là tiền tố (prefix) hay còn gọi là "Custom Variable Class". Có thể đặt là `my_app`, `security`, hay bất cứ thứ gì.
    + `current_tenant_id`: Là tên biến cụ thể muốn lưu.
    + nhãn này được lưu trong RAM
- Trong bảng `documents` có cột `tenant_id`. Thực hiện: 
```bash
-- 1. Bật bảo mật tầng dòng cho bảng
ALTER TABLE documents ENABLE ROW LEVEL SECURITY;

-- 2. Tạo một "Chính sách" (Policy)
CREATE POLICY tenant_policy ON documents
USING (tenant_id = current_setting('app.current_tenant_id'));
```
- Hàm `current_setting(...)` là lệnh của Postgres để đọc cái nhãn trong RAM đã dán ở bước SET LOCAL

2. Bọc Transaction cho phần bảo mật trước đó 
- WHY? Khi connection đó được gán nhãn chỉ dùng cho user được chỉ định, cái nhãn sẽ tồn tại cho đến khi connection biến mất. Tức là ngay cả khi connection trả về Pool thì cái nhãn vẫn không mất! Hậu quả khi user khác dùng nhưng cái nhãn vẫn còn.
- Vì vậy, khi bọc 1 Transaction, cái nhãn sẽ chỉ hiệu nghiệm trong Transaction đó và biến mất sau khi xong Transaction này. Khi trả về Pool, cái nhãn không còn và user khác có thể dùng nó