# RAG Evaluation Module (Go)

Module này cung cấp bộ công cụ đánh giá toàn diện cho hệ thống RAG (Retrieval-Augmented Generation), được thiết kế độc lập và linh hoạt để có thể sử dụng cho nhiều dự án khác nhau.

## 1. Cấu trúc dữ liệu chuẩn (Input)

Mọi bài đánh giá đều bắt đầu từ việc chuẩn bị danh sách `DatasetItem`:

```go
type DatasetItem struct {
    Question           string   // Câu hỏi của người dùng
    GroundTruthAnswer  string   // Câu trả lời mẫu chuẩn (để đối chiếu)
    GroundTruthContext string   // Văn bản gốc chuẩn (để đánh giá Retrieval)
    RetrievedContexts  []string // Danh sách văn bản tìm thấy từ hệ thống Search
    GeneratedAnswer    string   // Câu trả lời do chatbot sinh ra
}
```

---

## 2. Các chỉ số Retrieval (Tìm kiếm & Xếp hạng)

Các chỉ số này đánh giá khả năng của hệ thống trong việc tìm đúng tài liệu liên quan.

### a. Hit Rate@K (Tỉ lệ tìm thấy)
*   **Cách tính**: Tỉ lệ phần trăm các câu hỏi mà văn bản chuẩn (`GroundTruthContext`) xuất hiện trong danh sách tìm kiếm được (`RetrievedContexts`).
*   **Ý nghĩa**: Trả lời câu hỏi "Hệ thống có tìm thấy thông tin cần thiết không?".
*   **Công thức**: `(Số lần tìm thấy / Tổng số câu hỏi) * 100`

### b. Precision@1 (Độ chính xác tại vị trí số 1)
*   **Cách tính**: Tỉ lệ phần trăm các câu hỏi mà văn bản chuẩn nằm ngay tại vị trí kết quả đầu tiên.
*   **Ý nghĩa**: Đánh giá độ tin cậy tuyệt đối của kết quả đầu tiên.
*   **Công thức**: `(Số lần Rank 1 / Tổng số câu hỏi) * 100`

### c. MRR (Mean Reciprocal Rank - Xếp hạng trung bình đảo nghịch)
*   **Cách tính**: Tính trung bình cộng của các giá trị `1/Rank`. 
    *   Nếu tìm thấy ở vị trí 1, điểm là 1.
    *   Nếu vị trí 2, điểm là 0.5.
    *   Nếu vị trí 3, điểm là 0.33.
*   **Ý nghĩa**: Đây là chỉ số quan trọng nhất để đánh giá **khả năng xếp hạng**. MRR càng cao nghĩa là hệ thống càng đưa kết quả đúng lên phía trên.
*   **Công thức**: `(Σ 1/Rank) / Tổng số câu hỏi`

---

## 3. Các chỉ số Generation (LLM-as-a-Judge)

Sử dụng LLM để đánh giá ngữ nghĩa của câu trả lời.

### a. Faithfulness (Tính trung thực / Chống ảo giác)
*   **Cơ chế**: So sánh **Generated Answer** với **Retrieved Contexts**.
*   **Logic**: Kiểm tra xem mọi thông tin trong câu trả lời có bằng chứng trong ngữ cảnh hay không. 
*   **Mục tiêu**: Phát hiện "Hallucination" (LLM tự bịa thông tin).

### b. Answer Relevancy (Sự phù hợp)
*   **Cơ chế**: So sánh **Generated Answer** với **Question**.
*   **Logic**: Đánh giá xem câu trả lời có giải quyết đúng trọng tâm yêu cầu của người dùng không, có bị lan man hay thiếu ý không.

---

## 4. Hướng dẫn sử dụng

### Bước 1: Thu thập dữ liệu
```go
items := []evaluation.DatasetItem{
    {
        Question: "Giá dịch vụ?",
        GroundTruthContext: "Giá cắt tóc là 100k",
        RetrievedContexts: []string{"Dịch vụ uốn 200k", "Giá cắt tóc là 100k"},
    },
}
```

### Bước 2: Chạy đánh giá Retrieval
```go
result := evaluation.EvaluateRetrieval(items)
```

### Bước 3: Chạy đánh giá Ngữ nghĩa (LLM)
```go
judge := evaluation.NewJudge(llmClient)
score, _ := judge.ScoreFaithfulness(ctx, item.GeneratedAnswer, item.RetrievedContexts[0])
```

---

## 5. Lưu ý khi triển khai
*   **isMatch**: Hàm so sánh văn bản hiện tại đang dùng logic so sánh chuỗi (Substring/Prefix). Với các hệ thống phức tạp, có thể thay thế bằng hàm so sánh ID tài liệu để chính xác 100%.
*   **Judge Model**: Khuyên dùng các model mạnh như Llama 3-70B hoặc GPT-4 để làm Judge nhằm có kết quả khách quan nhất.
