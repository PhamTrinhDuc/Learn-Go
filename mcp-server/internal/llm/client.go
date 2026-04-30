package llm

import (
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
)

type Client struct {
	llm      llms.Model
	embedder embeddings.Embedder
}

func NewClient(cfg Config) (*Client, error) {
	// 1. Khởi tạo LLM
	llmModel, err := NewLLM(cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("llm init error: %w", err)
	}

	// 2. Khởi tạo Embedder
	embedder, err := NewEmbedder(cfg.Embed)
	if err != nil {
		return nil, fmt.Errorf("embedder init error: %w", err)
	}

	return &Client{
		llm:      llmModel,
		embedder: embedder,
	}, nil
}

