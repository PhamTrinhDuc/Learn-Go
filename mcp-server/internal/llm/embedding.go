package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	EMBEDDING_DIM = 1024
)

func NewEmbedder(cfg ProviderConfig) (embeddings.Embedder, error) {
	// 1. Sử dụng EmbedderClient làm interface trung gian
	var embedClient embeddings.EmbedderClient
	var err error

	switch cfg.Provider {
	case ProviderOllama:
		embedClient, err = ollama.New(
			ollama.WithServerURL(cfg.BaseURL),
			ollama.WithModel(cfg.Model),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init ollama embedding: %w", err)
		}

	case ProviderOpenAI, ProviderGroq:
		opts := []openai.Option{
			openai.WithToken(cfg.APIKey),
			openai.WithModel(cfg.Model),
			openai.WithEmbeddingDimensions(EMBEDDING_DIM),
		}
		if cfg.Provider == ProviderGroq {
			opts = append(opts, openai.WithBaseURL("https://api.groq.com/openai/v1"))
		} else if cfg.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
		}

		embedClient, err = openai.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to init openai/groq embedding: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", cfg.Provider)
	}

	// 2. Trả về trực tiếp vì NewEmbedder trả về (embeddings.Embedder, error)
	return embeddings.NewEmbedder(embedClient)
}

// GenerateEmbeddings chuyển văn bản thành vector
func (c *Client) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	return c.embedder.EmbedDocuments(ctx, texts)
}
