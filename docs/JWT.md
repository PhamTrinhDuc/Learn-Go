Luồng thực tế là:

1. Frontend POST /login với username/password
2. **Backend** dùng `privateKey` để **ký** và tạo ra token `xxx.yyy.zzz`
3. Backend trả token đó về cho frontend
4. Frontend lưu token vào bộ nhớ (localStorage, cookie, v.v.)
5. Mỗi request sau đó, frontend chỉ việc gắn vào header: `Authorization: Bearer xxx.yyy.zzz`

Frontend không hề biết `privateKey` là gì, cũng không tạo gì cả. Nó chỉ **cầm** cái token mà backend đã tạo sẵn.

---

**Vậy token `xxx.yyy.zzz` là gì?**

JWT gồm 3 phần nối nhau bằng dấu `.`:

- `xxx` = **Header** — thuật toán ký (ví dụ: RS256)
- `yyy` = **Payload** — dữ liệu: userID, tên, quyền, hạn dùng...
- `zzz` = **Signature** — chữ ký số, được tạo bằng `privateKey` của backend

Khi frontend gửi token lên, `JWTValidator` (tức cái Authentication Filter trong sơ đồ của bạn) sẽ dùng `publicKey` để **kiểm tra chữ ký** — nếu khớp thì biết token này do chính backend ký, không bị giả mạo.

---

**Tóm lại phân công rõ ràng:**

| Bên | Làm gì |
|---|---|
| Backend (`setupAuth`) | Tạo khóa, ký token, xác minh token |
| Frontend | Gửi login, nhận token, lưu lại, gắn vào mỗi request |

Frontend không cần biết gì về RSA hay secret key — đó chính là ưu điểm của JWT.