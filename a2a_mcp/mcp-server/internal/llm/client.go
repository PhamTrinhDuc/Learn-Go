package ollama

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Client struct {
	llm        *ollama.LLM
	embedder   *embeddings.EmbedderImpl
	llmModel   string
	embedModel string
}

func NewClient(cfg Config) (*Client, error) {
	// 1. Khởi tạo LLM
	llm, err := ollama.New(
		ollama.WithServerURL(cfg.BaseURL),
		ollama.WithModel(cfg.LLMModel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init llm: %w", err)
	}

	// 2. Khởi tạo Embedding model
	// Lưu ý: LangChainGo dùng chung struct ollama cho cả 2, nhưng ta sẽ bọc lại
	embedLLM, err := ollama.New(
		ollama.WithServerURL(cfg.BaseURL),
		ollama.WithModel(cfg.EmbedModel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init embedding model: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(embedLLM)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	return &Client{
		llm:        llm,
		embedder:   embedder,
		llmModel:   cfg.LLMModel,
		embedModel: cfg.EmbedModel,
	}, nil
}

// GenerateEmbeddings chuyển văn bản thành vector
func (c *Client) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	return c.embedder.EmbedDocuments(ctx, texts)
}

// Chat gọi LLM để trả lời câu hỏi (Blocking)
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx, c.llm, prompt,
		llms.WithTemperature(0), // Giảm sáng tạo để chính xác hơn
		// llms.WithRepeatPenalty(1.2), // Phạt nặng nếu lặp lại từ
	)
	if err != nil {
		return "", err
	}
	return completion, nil
}

// StreamChat gọi LLM và trả về kết quả từng từ một qua callback
func (c *Client) StreamChat(ctx context.Context, prompt string, callback func(string)) (string, error) {
	fullResp := ""
	_, err := c.llm.Call(ctx, prompt,
		llms.WithTemperature(0),
		// llms.WithRepeatPenalty(1.2),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			content := string(chunk)
			fullResp += content
			if callback != nil {
				callback(content)
			}
			return nil
		}),
	)
	return fullResp, err
}
