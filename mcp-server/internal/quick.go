package main

import (
	"context"
	"fmt"
	"log"
	"mcp-server/internal/database"
	llm "mcp-server/internal/llm"
	utils "mcp-server/internal/utils"
)

func getTestDBConfig() database.DBConfig {
	return database.DBConfig{
		Host:     utils.GetEnvString("DB_HOST", "localhost"),
		Port:     utils.GetEnvInt("DB_PORT", 5433),
		User:     utils.GetEnvString("DB_USER", "mcp_user"),
		Password: utils.GetEnvString("DB_PASSWORD", "mcp_password"),
		DBName:   utils.GetEnvString("DB_NAME", "salon_chain"),
		SSLMode:  utils.GetEnvString("DB_SSLMODE", "disable"),
		MaxConns: int32(utils.GetEnvInt("MAX_CONNS", 10)),
		MinConns: int32(utils.GetEnvInt("MAX_CONNS", 2)),
	}
}

func getTestModelConfig() llm.Config {
	return llm.Config{
		LLM: llm.ProviderConfig{
			Provider: llm.ProviderGroq,
			Model:    "llama-3.3-70b-versatile",
			APIKey:   utils.GetEnvString("GROQ_API_KEY", ""),
		},
		Embed: llm.ProviderConfig{
			Provider: llm.ProviderOllama,
			BaseURL:  "http://localhost:11434",
			Model:    "qwen3-embedding:0.6b",
		},
	}
}

func setupTest() (*database.DB, *llm.Client, error) {
	cfg := getTestDBConfig()
	db, err := database.NewDB(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init database")
	}
	model, err := llm.NewClient(getTestModelConfig())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init ollama model")
	}
	return db, model, nil
}

func TestPostgres(query string) {
	ctx := context.Background()

	db, _, err := setupTest()
	if err != nil {
		fmt.Println(err)
	}
	model, err := llm.NewClient(getTestModelConfig())
	if err != nil {
		fmt.Println("failed to init client: %w", err)
	}
	embeddings, err := model.GenerateEmbeddings(ctx, []string{query})
	if err != nil {
		fmt.Println("failed to embedding query: %w", err)
	}

	response, err := db.HybridSearch(
		ctx,
		database.HybridSearchParams{
			Query:        query,
			BM25Weight:   0.5,
			VectorWeight: 0.5,
			Embedding:    embeddings[0],
			Limit:        5,
		})

	// response, err := db.VectorSearch(ctx, embeddings[0], 5)

	if err != nil {
		fmt.Println("failed to perform vector search: %w", err)
	}

	// 1. Gom context từ kết quả tìm được
	contextText := ""
	for i, res := range response {
		contextText += fmt.Sprintf("\n--- Context %d ---\n%s\n", i+1, res.Document.Content)
		contextText += fmt.Sprintf("Reference %d %s\n", i+1, res.Document.Metadata["source"])
	}

	fmt.Println(contextText)

	// 2. Chuẩn bị Prompt hoàn chỉnh
	// fullPrompt := fmt.Sprintf(`
	// Bạn là trợ lý AI thông minh. Hãy trả lời câu hỏi của người dùng một cách chính xác dựa trên thông tin được cung cấp trong phần context.
	// Nếu thông tin trong context không đủ để trả lời, hãy nói rằng bạn không biết, đừng tự bịa ra câu trả lời.

	// CÂU HỎI: %s

	// CONTEXT:
	// %s
	// `, query, contextText)

	// // 3. Gọi LLM để lấy phản hồi (dùng StreamChat để thấy kết quả ngay lập tức)
	// fmt.Println("\n[*] Đang gửi dữ liệu cho LLM xử lý (Streaming)...")
	// fmt.Print("\n================= LLM RESPONSE =================\n")

	// aiResponse, err := model.StreamChat(ctx, fullPrompt, func(chunk string) {
	// 	fmt.Print(chunk) // In từng từ ra màn hình ngay khi AI sinh ra
	// })

	// fmt.Println(aiResponse)

	// if err != nil {
	// 	fmt.Printf("\nLỗi khi chat với LLM: %v\n", err)
	// 	return
	// }
	// fmt.Println("\n================================================")
}

func TestAIService() {
	ctx := context.Background()

	// 1. Khởi tạo client với cấu hình mặc định (qwen3.5 & qwen3-embedding)
	client, err := llm.NewClient(llm.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	// 2. Test thử sinh Embedding
	fmt.Println("[*] Đang thử sinh Embedding cho 1 câu văn...")
	vectors, err := client.GenerateEmbeddings(ctx, []string{"Chào bác, tôi là AI hỗ trợ Go!"})
	if err != nil {
		fmt.Printf("[!] Lỗi khi sinh embedding: %v\n", err)
	} else {
		fmt.Printf("[OK] Đã sinh xong vector, độ dài: %d\n", len(vectors[0]))
	}

	// 3. Test thử Chat LLM
	fmt.Println("[*] Đang hỏi Qwen2.5...")
	resp, err := client.Chat(ctx, "Chào bạn.")
	if err != nil {
		fmt.Printf("[!] Lỗi khi gọi LLM: %v\n", err)
	} else {
		fmt.Printf("[Qwen]: %s\n", resp)
	}
}

func main() {
	TestPostgres("Chăm sóc tóc layer thế nào?")
}
