# Postgres Hybrid Search Architecture (Deep Dive)

Tài liệu này không chỉ hướng dẫn cách dùng, mà còn mổ xẻ tận gốc rễ (root cause) thiết kế kiến trúc bên dưới hệ thống tìm kiếm Hybrid của PostgreSQL, loại bỏ hoàn toàn các khái niệm "hộp đen" (blackbox).

## 1. Khởi nguồn: Bài toán Full-Text Search (FTS)

Khi tìm kiếm chuỗi thông thường (dùng `LIKE` hoặc `ILIKE`), hệ quản trị CSDL sẽ phải duyệt qua từng dòng một (Sequential Scan), dò từng ký tự. Điều này vừa chậm (O(N) thay vì O(1) hoặc O(logN)), vừa ngu ngốc (vì người dùng gõ "mèo" hay "những con mèo" thì LIKE không hiểu được đó là cùng một thực thể).

Giải pháp là FTS: Phân tích ngôn ngữ trước khi lưu trữ (Tokenization & Lexing) và tạo mục lục ngược (Inverted Index).

---

## 2. Giải phẫu quá trình "Nén" văn bản (`to_tsvector`)

Khi bạn gọi `to_tsvector('Cấu_hình', 'Nội dung văn bản')`, chuỗi văn bản không được lưu giấu đi một cách ma thuật nào cả. Nó đi qua một **Đường ống (Pipeline) 3 bước khắt khe**:

### Bước A: The Parser (Bộ phân chia)
Văn bản thô đi vào Parser. Parser không quan tâm ý nghĩa, nó chỉ dùng dấu phân cách (khoảng trắng, dấu câu) và nhận diện các loại chuỗi (số, email, url, chữ).
*   Đầu vào: `"Tạm ứng lương 2024!"`
*   Đầu ra: Token 1 (Chữ): `Tạm`, Token 2 (Chữ): `ứng`, Token 3 (Chữ): `lương`, Token 4 (Số): `2024`.

### Bước B: Dictionaries & Stemming (Bộ từ điển & Trích xuất gốc từ) - Tại sao cấu hình 'english' lại "hại" tiếng Việt?

Đây là phần "não" của FTS.
*   **Stemming là gì?** Stemming là quá trình ngôn ngữ học máy tính nhằm cắt bỏ các phần phụ tố (tiền tố/hậu tố) để đưa một từ về gốc (Root/Stem). Ví dụ tiếng Anh: `running`, `runs`, `ran` -> gốc từ là `run`. Điều này giúp người dùng search "run" vẫn ra văn bản có chứa "running".
*   **Cấu hình `'english'` hoạt động thế nào?** Nó áp dụng thuật toán stemming Snowball cho chuẩn tiếng Anh. Khi gặp từ `thưởng` (tiếng Việt), thuật toán Snowball không hiểu, nó cố gắng gọt giũa một cách mù quáng (ví dụ cắt mất chữ `g` hoặc biến dạng từ), dẫn đến Token sinh ra chứa rác hoặc sai lệch.
*   **Tại sao dùng `'simple'` lại cứu được tiếng Việt?** Cấu hình `'simple'` có một bộ Dictionary rất đơn giản: **Không làm gì cả**. Nó cấm Stemming. Nó chỉ đơn thuần lấy Token từ Bước A, chuyển tất cả thành chữ thường (lowercase) và ném ra ngoài. Đối với ngôn ngữ đơn lập, không có biến đổi hình thức từ như tiếng Việt, việc *không làm gì* hóa ra lại là cách xử lý đúng đắn nhất.

### Bước C: Tạo Lexemes (Cấu trúc dữ liệu cuối cùng)
Lexemes là các từ đã qua xử lý ở Bước B, kèm theo vị trí của chúng.
Đầu ra cuối cùng của `to_tsvector('simple', 'Tạm ứng lương')` là:
`'luong':3 'tam':1 'ung':2` (Giả sử ta đã lột dấu).

---

## 3. Giải phẫu `unaccent`: Cơ chế lột dấu

Bạn tự hỏi hàm `unaccent` (hay `f_unaccent` wrapper chúng ta viết) hoạt động thế nào ở mức độ máy tính?

*   Nó không dùng AI để hiểu tiếng Việt. Nó hoàn toàn là **Mapping ký tự cứng (Character translation rules)**.
*   Bên trong mã nguồn Postgres extension, có một file ruleset (vd: `unaccent.rules`). Trong đó định nghĩa rõ: Byte chứa ký tự `đ` -> chuyển thành `d`; Byte chứa ký tự `ự` -> chuyển thành `u`.
*   Khi chạy qua `unaccent`, chuỗi `Tạm ứng` sẽ bị quét từng ký tự, đối chiếu với bảng ruleset trong RAM (bảng hash) và mapping sang `Tam ung`. Thao tác này là O(N) với N là độ dài văn bản, rất nhẹ và tĩnh (đó là lý do ta có thể gói nó vào hàm `RÉTURN IMMUTABLE` để tạo Index được).

---

## 4. Xếp hạng (Ranking): Thuật toán `ts_rank_cd`

Tại sao ta dùng `ts_rank_cd` thay vì `ts_rank` thường? `cd` là viết tắt của **Cover Density (Mật độ che phủ)**.

### Toán học đằng sau Cover Density
`ts_rank` thông thường đếm số lượng từ vựng xuất hiện (Tần suất - TF).
`ts_rank_cd` thì thông minh hơn: Nó đo lường **KHOẢNG CÁCH** giữa các từ khóa (Cover).
*   Giả sử user search: `'quy' & 'định' & 'thưởng'`
*   Tài liệu 1: `"Quy định về việc nhân viên được thưởng..."` => Các từ trải dài cách nhau 4-5 từ. Khoảng cách (Cover) lớn => Điểm thấp.
*   Tài liệu 2: `"Hôm nay ra quy định thưởng..."` => Các từ khoá nằm sát rạt nhau (Cover hẹp lại, Density dày lên). Mật độ dày đặc => Điểm cực cao.

### Tại sao nó mạnh khi tìm Mã (ID), Mã tài liệu, Term kỹ thuật?
Khi người dùng gõ tìm mã `QyĐ-VTI-01`, qua quá trình parsing, nó biến thành các cụm tokens đứng sát nhau. `ts_rank_cd` ưu tiên tuyệt đối các block text nơi các token này nén chặt lại với nhau như một cụm từ gốc, thay vì các token rải rác. Do vậy tính chính xác (Precision) cực kỳ cao (gần như Exact Match), bù đắp hoàn hảo cho Vector Search (thường hay bị ảo giác, tìm liên miên vạn dặm).

---

## 5. Kiến trúc Inverted Index (GIN) = Chìa khoá của tốc độ

Bạn gắn `USING gin` vào lúc tạo Index. GIN (Generalized Inverted Index) là gì?
Nó giống như **bảng Mục Lục** ở cuối quyển bách khoa toàn thư.
Thay vì lưu: `Dòng 1: [lương, thưởng, tiền]`
GIN lưu ngược lại:
*   `lương -> Dòng 1, Dòng 80, Dòng 100`
*   `thưởng -> Dòng 1, Dòng 5`

Khi user search `'lương' & 'thưởng'`, Postgres vào thẳng bảng GIN, nhặt ra list của `lương` và list của `thưởng`, thực hiện phép toán giao (Intersection) giữa 2 mảng bit. Pằng! Dòng 1 văng ra ngay lập tức với tốc độ O(1) tra cứu cây, bất chấp CSDL có hàng chục triệu row.

---

## 6. Tổng hợp: Hybrid Search (RRF - Reciprocal Rank Fusion)

Khi bạn đã hiểu sâu FTS, bạn nhận ra: FTS vô địch ở mức độ Token (Lexical). Nhưng nếu user gõ `"tiền bo"` tài liệu có chữ `"thưởng"` sẽ tạch, vì 2 token này khác nhau byte-mapping.

Vì vậy ta cần nhúng Vector Search bằng AI (Ollama). Vector Search đo góc Cosine giữa 2 mảng ma trận đa chiều, nghĩa là từ góc độ toán học đại số, mảng vector của từ "tiền bo" và "thưởng" nằm rất sát nhau (tính toán qua Attention mechanism trong Neural Network). 

**Vấn đề:** 
* DBMS FTS trả về điểm BM25 Scale từ 0 -> 1000+.
* Qwen Embedding trả về Cosine Distance Scale từ 0 -> 1.
Chúng không cùng hệ quy chiếu, cộng lại là thảm hoạ toán học.

**Giải pháp cốt lõi: RRF**
RRF dẹp bỏ "Điểm tuyệt đối", chỉ dùng "Bảng xếp hạng" (Thứ hạng 1, 2, 3).
Công thức chuẩn mực: `Score_RRF(item) = 1 / (K + Rank_BM25(item)) + 1 / (K + Rank_Vector(item))` (thường K = 60).
Bằng cách đẩy mọi thứ xuống cùng một hệ quy chiếu "Xếp hạng", ta kết hợp được sức mạnh cơ bắp tra cứu cấu trúc tĩnh tuyệt đối của PG GIN với sức mạnh tư duy mở đại số của LLMs.
