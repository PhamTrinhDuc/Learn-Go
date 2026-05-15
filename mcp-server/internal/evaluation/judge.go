package evaluation

import (
	"context"
	"encoding/json"
	"fmt"
	"mcp-server/internal/llm"
)

// GenerationMetrics lưu trữ kết quả đánh giá phần sinh nội dung bằng LLM
type GenerationMetrics struct {
	Faithfulness    float64 `json:"faithfulness"`
	AnswerRelevancy float64 `json:"answer_relevancy"`
}

// Judge chịu trách nhiệm gọi LLM để chấm điểm các chỉ số ngữ nghĩa
type Judge struct {
	client *llm.OpenAICompatibleLLM
}

func NewJudge(client *llm.OpenAICompatibleLLM) *Judge {
	return &Judge{client: client}
}

// ScoreFaithfulness đánh giá xem câu trả lời có dựa trên ngữ cảnh hay không (Hallucination check)
func (j *Judge) ScoreFaithfulness(ctx context.Context, answer, retrievedContext string) (float64, error) {
	prompt := fmt.Sprintf(`
BẠN LÀ MỘT GIÁM KHẢO ĐÁNH GIÁ HỆ THỐNG RAG. 
NHIỆM VỤ: Kiểm tra tính trung thực (Faithfulness) của câu trả lời dựa trên ngữ cảnh được cung cấp.

[Context]: %s
[Answer]: %s

YÊU CẦU:
1. Chỉ chấm điểm 1.0 nếu TẤT CẢ các ý chính trong Câu trả lời đều có thể tìm thấy hoặc suy luận trực tiếp từ Ngữ cảnh.
2. Nếu có bất kỳ thông tin nào bị bịa đặt hoặc không có trong Ngữ cảnh, hãy trừ điểm.
3. Trả về kết quả dưới định dạng JSON duy nhất: {"score": 0.0 - 1.0}
`, retrievedContext, answer)

	return j.getScore(ctx, prompt)
}

// ScoreAnswerRelevancy đánh giá xem câu trả lời có thực sự giải quyết được câu hỏi không
func (j *Judge) ScoreAnswerRelevancy(ctx context.Context, question, answer string) (float64, error) {
	prompt := fmt.Sprintf(`
BẠN LÀ MỘT GIÁM KHẢO ĐÁNH GIÁ HỆ THỐNG RAG. 
NHIỆM VỤ: Đánh giá mức độ phù hợp (Relevancy) của câu trả lời đối với câu hỏi của người dùng.

[Question]: %s
[Answer]: %s

YÊU CẦU:
1. Điểm cao (1.0) nếu câu trả lời đi thẳng vào vấn đề và đầy đủ thông tin khách cần.
2. Điểm thấp nếu câu trả lời lan man hoặc không giải quyết được nhu cầu của câu hỏi.
3. Trả về kết quả dưới định dạng JSON duy nhất: {"score": 0.0 - 1.0}
`, question, answer)

	return j.getScore(ctx, prompt)
}

// Helper để parse kết quả từ LLM
func (j *Judge) getScore(ctx context.Context, prompt string) (float64, error) {
	resp, err := j.client.Chat(ctx, prompt)
	if err != nil {
		return 0, err
	}

	var result struct {
		Score float64 `json:"score"`
	}
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return 0, fmt.Errorf("failed to parse judge score: %v", err)
	}

	return result.Score, nil
}
