package llm

import (
	"context"
	"fmt"
	"log"
	"mcp-server/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAIService(t *testing.T) {
	ctx := context.Background()

	// 1. Khởi tạo client với cấu hình mặc định (qwen2.5 & qwen3-embedding)
	client, err := NewClient(Config{
		LLM: ProviderConfig{
			Provider: ProviderGroq,
			Model:    "llama-3.3-70b-versatile",
			APIKey:   utils.GetEnvString("GROQ_API_KEY", ""),
		},
		Embed: ProviderConfig{
			Provider: ProviderOllama,
			BaseURL:  "http://localhost:11434",
			Model:    "qwen3-embedding:0.6b",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Test thử sinh Embedding
	fmt.Println("[*] Đang thử sinh Embedding cho 1 câu văn...")
	vectors, err := client.GenerateEmbeddings(ctx, []string{"Chào bác, tôi là AI hỗ trợ Go!"})
	assert.NoError(t, err)
	assert.Condition(t, func() (success bool) {
		return len(vectors) > 0
	}, "Not found vector embedding")

	// 3. Test thử Chat LLM
	fmt.Println("[*] Đang hỏi LLM...")
	resp, err := client.Chat(ctx, "Chào bác, bác làm được gì?")
	fmt.Println(resp)
	assert.NoError(t, err)
	assert.Condition(t, func() (success bool) {
		return len(resp) > 0
	}, "Response should be have len > 0")
}
