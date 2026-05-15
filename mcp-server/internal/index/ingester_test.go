package index

import (
	"context"
	database "mcp-server/internal/database"
	llm "mcp-server/internal/llm"
	"mcp-server/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	db, err := database.NewDB(ctx, database.NewDBConfig())
	assert.NoError(t, err)
	// 2. Init client Ollama
	model := llm.NewOpenAICompatibleEmbedding(llm.NewOpenAIEmbeddingConfig())
	assert.NoError(t, err)

	err = ingestBatchingToVDB(ctx, model, db, filePaths)
	assert.NoError(t, err)
}
