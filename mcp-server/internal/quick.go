package main

import (
	"context"
	"fmt"
	database "mcp-server/internal/database"
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
}

func TestBenchmark() {
	filePath := "../data/evals/eval_dataset.csv"
	// err := database.GenDataset(filePath)
	// if err != nil {
	// 	fmt.Printf("failed to gen dataset: %s", err)
	// }

	err := database.Evaluation(filePath, true)
	if err != nil {
		fmt.Println("failed to evaluaton dataset: %w", err)
	}
}

func main() {
	// TestPostgres("Chăm sóc tóc layer thế nào?")
	TestBenchmark()
}
