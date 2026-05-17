package evaluation

import (
	"context"
	"encoding/json"
	"fmt"
	"mcp-server/internal/llm"
	"strings"
)

// GenerationMetrics lưu trữ kết quả đánh giá phần sinh nội dung bằng LLM
type GenerationMetrics struct {
	Faithfulness    float64 `json:"faithfulness"`
	AnswerRelevancy float64 `json:"answer_relevancy"`
}

// Judge chịu trách nhiệm gọi LLM để chấm điểm các chỉ số ngữ nghĩa
type Judge struct {
	client llm.LLMModel
}

func NewJudge(client llm.LLMModel) *Judge {
	return &Judge{client: client}
}

// ScoreFaithfulness đánh giá xem câu trả lời có dựa trên ngữ cảnh hay không (Hallucination check)
func (j *Judge) ScoreFaithfulness(ctx context.Context, answer string, retrievedContext []string) (float64, error) {
	prompt := fmt.Sprintf(`
You are a judge evaluating the faithfulness of an answer to a user's question.
Mission: Rate the answer to question.

[Context]: %s
[Answer]: %s

### REQUIRE IMPORTANT:
1. The answer must directly address the user's question.
2. The answer must be comprehensive and complete.
3. The answer must be accurate and not contain any false information.
4. The answer must be concise and easy to understand.

Return result in JSON format: 
{{
	"score": 0.0 - 1.0
}}
`, retrievedContext, answer)

	return j.getScore(ctx, prompt)
}

// ScoreAnswerRelevancy đánh giá xem câu trả lời có thực sự giải quyết được câu hỏi không
func (j *Judge) ScoreAnswerRelevancy(ctx context.Context, question, answer string) (float64, error) {
	prompt := fmt.Sprintf(`
You are a judge evaluating the relevancy of an answer to a user's question.
Mission: Rate the answer to question.

[Question]: %s
[Answer]: %s

### REQUIRE IMPORTANT:
1. The answer must directly address the user's question.
2. The answer must be comprehensive and complete.
3. The answer must be accurate and not contain any false information.
4. The answer must be concise and easy to understand.

Return result in JSON format: 
{{
	"score": 0.0 - 1.0
}}
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
	scoreStr := strings.TrimSpace(resp)
	scoreStr = strings.TrimPrefix(scoreStr, "```json")
	scoreStr = strings.TrimSuffix(scoreStr, "```")
	if err := json.Unmarshal([]byte(scoreStr), &result); err != nil {
		return 0, fmt.Errorf("failed to parse judge score: %v", err)
	}

	return result.Score, nil
}
