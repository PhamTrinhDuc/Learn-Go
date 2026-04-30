package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

func NewLLM(cfg ProviderConfig) (llms.Model, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("missing API Key for provider: %s", cfg.Provider)
	}
	switch cfg.Provider {
	case ProviderOllama:
		return ollama.New(
			ollama.WithServerURL(cfg.BaseURL),
			ollama.WithModel(cfg.Model),
		)
	case ProviderOpenAI, ProviderGroq:
		opts := []openai.Option{
			openai.WithToken(cfg.APIKey),
			openai.WithModel(cfg.Model),
		}
		if cfg.Provider == ProviderGroq {
			opts = append(opts, openai.WithBaseURL("https://api.groq.com/openai/v1"))
		} else if cfg.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
		}
		return openai.New(opts...)
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", cfg.Provider)
	}
}

// Chat gọi LLM để trả lời câu hỏi (Blocking)
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx, c.llm, prompt,
		llms.WithTemperature(0),
	)
	if err != nil {
		return "", err
	}
	return completion, nil
}

// StreamChat gọi LLM và trả lời câu hỏi (Streaming)
func (c *Client) StreamChat(ctx context.Context, prompt []llms.MessageContent, callback func(string)) (string, error) {
	fullResp := ""
	_, err := c.llm.GenerateContent(ctx, prompt,
		llms.WithTemperature(0),
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
