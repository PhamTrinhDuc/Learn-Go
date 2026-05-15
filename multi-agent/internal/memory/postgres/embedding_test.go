package postgres

import (
	"context"
	"multi-agent/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbed(t *testing.T) {
	embeddingDim := 1024
	embeder := NewOpenAICompatibleEmbedding(
		OpenAICompatibleEmbeddingConfig{
			// BaseURL:   utils.GetEnvString("BASE_URL_OLLAMA", "http://localhost:11434/v1"),
			// Model:     utils.GetEnvString("EMBEDDING_MODEL", "qwen3-embedding:0.6b"),
			BaseURL:   utils.GetEnvString("BASE_URL_OPENAI", "https://api.openai.com/v1"),
			Model:     "text-embedding-3-large",
			Dimension: utils.GetEnvInt("EMBEDDING_DIM", embeddingDim),
			APIKey:    utils.GetEnvString("OPENAI_API_KEY", ""),
		},
	)

	embedding, err := embeder.Embed(context.Background(), "Hello")
	assert.NoError(t, err)
	assert.Equal(t, len(embedding), embeddingDim)
}
