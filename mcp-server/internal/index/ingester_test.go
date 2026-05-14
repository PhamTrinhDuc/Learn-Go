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

var configModel = llm.Config{
	LLM: llm.ProviderConfig{
		Provider: llm.ProviderGroq,
		Model:    "llama-3.3-70b-versatile",
		APIKey:   utils.GetEnvString("GROQ_API_KEY", ""),
	},
	Embed: llm.ProviderConfig{
		Provider: llm.ProviderOpenAI,
		// BaseURL:  "http://localhost:11434",
		Model:  "text-embedding-3-small",
		APIKey: utils.GetEnvString("OPENAI_API_KEY", ""),
	},
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
	model, err := llm.NewClient(configModel)
	assert.NoError(t, err)

	err = ingestBatchingToVDB(ctx, model, db, filePaths)
	assert.NoError(t, err)
}
