package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbed(t *testing.T) {
	embeddingDim := 1024
	embeder := NewOpenAICompatibleEmbedding(
		OpenAICompatibleEmbeddingConfig{
			BaseURL:   "http://localhost:11434/v1",
			Model:     "qwen3-embedding:0.6b",
			Dimension: embeddingDim,
		},
	)

	embedding, err := embeder.Embed(context.Background(), "Hello")
	assert.NoError(t, err)
	assert.Equal(t, len(embedding), embeddingDim)
}
