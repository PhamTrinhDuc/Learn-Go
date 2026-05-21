# HỌC VIỆN CÔNG NGHỆ BƯU CHÍNH VIỄN THÔNG
## KHOA CÔNG NGHỆ THÔNG TIN 1

**CỘNG HOÀ XÃ HỘI CHỦ NGHĨA VIỆT NAM**  
**Độc lập – Tự do – Hạnh phúc**

---

# ĐỀ CƯƠNG NGHIÊN CỨU KHOA HỌC SINH VIÊN 2025

---

**1. Tên đề tài:**  
**Chưng cất tri thức từ mô hình lai CNN-Mamba sang mạng CNN nhẹ cho bài toán phân đoạn ảnh y tế**  
*(Knowledge Distillation from CNN-Mamba Hybrid to Lightweight CNN for Medical Image Segmentation)*

**Mã số:** ……… – SV – 2025 – CNTT1

---

**2. Thuộc lĩnh vực:** RD – Đào tạo

---

**3. Mục tiêu, nội dung và kết quả đề tài:**

**a. Mục tiêu:**

Nghiên cứu và triển khai kỹ thuật chưng cất tri thức (Knowledge Distillation – KD) từ mô hình teacher lai CNN-Mamba (CNNMamba) sang mô hình student CNN nhẹ (LightCNN) cho bài toán phân đoạn ảnh y tế, nhằm đạt được sự cân bằng giữa hiệu năng phân đoạn và chi phí tính toán. Đề tài hướng đến xây dựng một mô hình student nhỏ gọn (~1.2 triệu tham số) có thể triển khai trong các môi trường tài nguyên hạn chế, đồng thời duy trì chất lượng phân đoạn gần với mô hình teacher (~4–5 triệu tham số).

**b. Nội dung:**

- Tìm hiểu kiến trúc mô hình Mamba (State Space Model) và mô hình lai CNN-Mamba trong bài toán phân đoạn ảnh y tế
- Tìm hiểu kỹ thuật chưng cất tri thức (Knowledge Distillation), bao gồm logit-based KD và feature-based KD
- Thiết kế và cài đặt mô hình teacher CNNMamba và mô hình student LightCNN
- Xây dựng hàm loss kết hợp: Dice Loss + Focal Loss + Boundary Loss cho bài toán phân đoạn với dữ liệu mất cân bằng lớp (class imbalance)
- Thiết kế pipeline chưng cất tri thức kết hợp logit KD (soft targets) và feature KD (feature map alignment)
- Thực hiện ablation study đánh giá ảnh hưởng của các siêu tham số KD: hệ số cân bằng alpha, nhiệt độ temperature, trọng số feature KD
- Đánh giá mô hình trên tập dữ liệu BRISC 2025 với các chỉ số: Dice, IoU, Precision, Recall, HD95
- So sánh hiệu năng giữa: Teacher (CNNMamba), Student + KD, Student không KD
- Viết báo cáo khoa học và bài báo tạp chí

**c. Kết quả:**

- Báo cáo khoa học đầy đủ
- Mã nguồn chương trình huấn luyện và đánh giá mô hình
- Bài báo gửi đăng tạp chí *Pattern Recognition and Image Analysis* (Springer, SCOPUS)

---

**4. Đơn vị chủ trì:** Khoa Công nghệ thông tin 1

**5. Đơn vị phối hợp** (nếu có): *(để trống)*

**6. Chủ trì đề tài:** *[Họ và tên bạn] – Mã sinh viên: [Mã SV] – Lớp: [Lớp]*

**7. Những người tham gia thực hiện:** *(để trống hoặc bổ sung nếu có)*

**8. Người hướng dẫn:** *[Tên cô hướng dẫn]*

---

**9. Sơ lược tình hình nghiên cứu trong nước, ngoài nước:**

**a. Sơ lược tình hình**

*Trong nước:*  
Các nghiên cứu về phân đoạn ảnh y tế tại Việt Nam chủ yếu tập trung vào ứng dụng các kiến trúc UNet và Transformer có sẵn. Chưa có nhiều công trình nghiên cứu kết hợp State Space Model (Mamba) với kỹ thuật chưng cất tri thức trong bối cảnh phân đoạn ảnh y tế, đặc biệt hướng đến mô hình nhẹ phù hợp triển khai thực tế trong môi trường y tế tài nguyên hạn chế.

*Ngoài nước:*  
Kiến trúc Mamba dựa trên State Space Model (SSM) gần đây đã cho thấy hiệu quả vượt trội trong xử lý chuỗi dài nhờ khả năng nắm bắt phụ thuộc tầm xa với độ phức tạp tuyến tính. Các mô hình như VMamba, U-Mamba, SegMamba đã được áp dụng trong phân đoạn ảnh y tế và đạt kết quả tốt. Song song đó, kỹ thuật chưng cất tri thức (Hinton et al., 2015) đã được nghiên cứu rộng rãi trong nén mô hình, tuy nhiên việc kết hợp KD từ mô hình lai CNN-Mamba sang CNN thuần vẫn còn ít được khai thác, đặc biệt trong lĩnh vực phân đoạn ảnh y tế với class imbalance nặng.

*Tài liệu tham khảo:*

- Hinton, G., et al. (2015). Distilling the Knowledge in a Neural Network. *NIPS Workshop*.
- Gu, A., & Dao, T. (2023). Mamba: Linear-Time Sequence Modeling with Selective State Spaces. *arXiv:2312.00752*.
- Liu, Y., et al. (2024). VMamba: Visual State Space Model. *NeurIPS 2024*.
- Ronneberger, O., et al. (2015). U-Net: Convolutional Networks for Biomedical Image Segmentation. *MICCAI 2015*.
- Lin, T.-Y., et al. (2017). Focal Loss for Dense Object Detection. *ICCV 2017*.

**b. Nhu cầu thực tế và khả năng áp dụng kết quả đề tài:**

- Phân đoạn khối u trong ảnh y tế là bước quan trọng trong hỗ trợ chẩn đoán lâm sàng, yêu cầu độ chính xác cao nhưng cũng cần có thể triển khai trên thiết bị có tài nguyên hạn chế
- Mô hình student sau KD (~1.2M tham số) có tiềm năng triển khai thực tế trong các hệ thống hỗ trợ chẩn đoán tại cơ sở y tế không có GPU cao cấp
- Kết quả đề tài đóng góp vào hướng nghiên cứu model compression trong y tế, một lĩnh vực đang được quan tâm ngày càng nhiều tại Việt Nam và quốc tế

---

**10. Sản phẩm giao nộp của đề tài:**

**Dạng sản phẩm III:** Chương trình máy tính; Phương pháp, giải pháp

**a. Kết quả, sản phẩm:**

| STT | Tên sản phẩm | Đơn vị tính | Số lượng | Chú thích |
|-----|--------------|-------------|----------|-----------|
| 1 | Mã nguồn pipeline huấn luyện (Teacher + Student + KD) | Bộ | 1 | Notebook Kaggle/Python |
| 2 | Mô hình đã huấn luyện (Teacher checkpoint, Student checkpoint) | File | 2 | .pth |
| 3 | Bài báo tạp chí khoa học | Bài | 1 | Gửi đăng Springer PRIA |
| 4 | Báo cáo kết quả đề tài | Quyển | 05 | |

**b. Yêu cầu khoa học, kinh tế xã hội đối với sản phẩm (dạng II, III):**

| STT | Tên sản phẩm | Yêu cầu cụ thể cần đạt | Chú thích |
|-----|--------------|------------------------|-----------|
| 1 | Mã nguồn | Tái hiện được kết quả, có tài liệu mô tả rõ ràng | |
| 2 | Mô hình Student + KD | Dice ≥ 0.83, HD95 ≤ 22, nhỏ hơn teacher ≥ 3 lần về tham số | |
| 3 | Bài báo tạp chí | Được chấp nhận đăng hoặc đang trong quá trình review tại PRIA – Springer | SCOPUS |
| 4 | Báo cáo kết quả đề tài | Khoa học, đầy đủ nội dung | |

---

**11. Tiến độ thực hiện các nội dung:**

| STT | Nội dung từng bước | Thời gian thực hiện | Kết quả cần đạt |
|-----|--------------------|---------------------|-----------------|
| 1 | Xây dựng và xét duyệt đề cương | 03–05/2025 | Đề cương được duyệt |
| 2 | Tìm hiểu lý thuyết: Mamba, KD, segmentation loss | 05/2025 | Tài liệu tổng quan |
| 3 | Cài đặt mô hình Teacher (CNNMamba), huấn luyện và đánh giá | 05–06/2025 | Teacher checkpoint, kết quả baseline |
| 4 | Cài đặt pipeline KD, ablation study các siêu tham số | 06–07/2025 | Bảng ablation, Student + KD tốt nhất |
| 5 | Đánh giá so sánh, visualize kết quả, viết paper | 07–08/2025 | Draft bài báo hoàn chỉnh |
| 6 | Submit bài báo, hoàn thiện báo cáo đề tài | 08/2025 | Bài báo đã submit, báo cáo hoàn thiện |
| 7 | Hội thảo / nghiệm thu đề tài | 11/2025 | Thông qua Hội đồng nghiệm thu |

---

**12. Dự toán kinh phí theo các khoản chi:**

| STT | Nội dung các khoản chi | Thành tiền |
|-----|------------------------|------------|
| 1 | Xây dựng và duyệt đề cương nghiên cứu | 195.000 |
| 2 | Viết báo cáo kết quả nghiên cứu | 850.000 |
| 3 | Hội thảo khoa học | 0 |
| 4 | Nghiệm thu | 455.000 |
| 5 | Chi phí văn phòng phẩm, photocopy, đóng quyển | 0 |
| **Tổng cộng** | | **1.500.000** |

---

**CHỦ TRÌ ĐỀ TÀI**  
*(Ký tên và ghi rõ họ tên)*

&nbsp;

**GIÁO VIÊN HƯỚNG DẪN**  
*(Ký tên và ghi rõ họ tên)*

&nbsp;

**ĐƠN VỊ CHỦ TRÌ ĐỀ TÀI**  
*Khoa Công nghệ thông tin 1*