package index

import (
	"context"
	database "mcp-server/internal/database"
	ollama "mcp-server/internal/llm"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configDB = database.DBConfig{
	Host:     "localhost",
	Port:     5433,
	User:     "mcp_user",
	Password: "mcp_password",
	DBName:   "mcp_db",
}

var configModel = ollama.Config{
	BaseURL:    "http://localhost:11434",
	LLMModel:   "qwen2.5:0.5b",
	EmbedModel: "qwen3-embedding:0.6b",
}

func TestIngestToVDB(t *testing.T) {

	ctx := context.Background()
	tenantID := "11111111-1111-1111-1111-111111111111"
	filePath := "../../../data/VTI_Quy định_Quy định điền thông tin trên hệ thống VMS_v2.0.pdf"

	// 1. Init database
	db, err := database.NewDB(ctx, configDB)
	assert.NoError(t, err)
	// 2. Init client Ollama
	model, err := ollama.NewClient(configModel)

	err = ingestToVDB(ctx, model, db, filePath, tenantID)
	assert.NoError(t, err)
}

func TestIngestBatchingToVDB(t *testing.T) {

	ctx := context.Background()
	tenantID := "11111111-1111-1111-1111-111111111111"
	filePaths := []string{
		"../../../data/VTI_Quy định thưởng đề xuất IP Kaizen-2019_V1.0.pdf",
		"../../../data/VTI_Quy định_Quy định điền thông tin trên hệ thống VMS_v2.0.pdf",
		"../../../data/VTI_Quy định_Quy định tạm ứng lương_v3.0.pdf",
		// "../../../data/VTI_Thỏa thuận tạm ứng lương_v1.0.pdf",
	}

	// 1. Init database
	db, err := database.NewDB(ctx, configDB)
	assert.NoError(t, err)
	// 2. Init client Ollama
	model, err := ollama.NewClient(configModel)

	err = ingestBatchingToVDB(ctx, model, db, filePaths, tenantID)
	assert.NoError(t, err)
}
