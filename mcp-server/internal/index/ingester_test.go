package index

import (
	"context"
	database "mcp-server/internal/database"
	llm "mcp-server/internal/llm"
	"mcp-server/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configDB = database.DBConfig{
	Host:     "localhost",
	Port:     5433,
	User:     "mcp_user",
	Password: "mcp_password",
	DBName:   "salon_chain",
}

var llmConfig = llm.OpenAICompatibleConfig{
	// BaseURL:   utils.GetEnvString("BASE_URL_OLLAMA", "http://localhost:11434/v1"),
	// Model:     utils.GetEnvString("EMBEDDING_MODEL", "qwen3-embedding:0.6b"),
	BaseURL:   utils.GetEnvString("BASE_URL_OPENAI", "https://api.openai.com/v1"),
	Model:     "text-embedding-3-large",
	Dimension: utils.GetEnvInt("EMBEDDING_DIM", 1024),
	APIKey:    utils.GetEnvString("OPENAI_API_KEY", "abcd"),
}

// func TestIngestToVDB(t *testing.T) {

// 	ctx := context.Background()
// 	filePath := "../../data/agent-instructions/booking-agent/quy-trinh-dat-lich.md"

// 	// 1. Init database
// 	db, err := database.NewDB(ctx, configDB)
// 	assert.NoError(t, err)
// 	// 2. Init client Ollama
// 	model, err := llm.NewClient(configModel)

// 	err = ingestToVDB(ctx, model, db, filePath)
// 	assert.NoError(t, err)
// }

func TestIngestBatchingToVDB(t *testing.T) {

	ctx := context.Background()
	dataRoot := "../../data"
	filePaths, err := utils.GetListFiles(dataRoot)
	assert.NoError(t, err)
	// 1. Init database
	db, err := database.NewDB(ctx, configDB)
	assert.NoError(t, err)
	// 2. Init client Ollama
	model := llm.NewOpenAICompatibleEmbedding(llmConfig)
	assert.NoError(t, err)

	err = ingestBatchingToVDB(ctx, model, db, filePaths)
	assert.NoError(t, err)
}
