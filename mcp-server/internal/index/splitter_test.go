package index

import (
	"context"
	"fmt"
	"mcp-server/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDocument(t *testing.T) {
	filePath := "../../data/agent-instructions/booking-agent/quy-trinh-dat-lich.md"
	docs, err := loadDocs(context.Background(), filePath)

	assert.NoError(t, err)
	assert.Condition(t, func() (success bool) {
		return len(docs) > 0
	}, "The len documents should be > 0")
}

func TestFormatDocument(t *testing.T) {
	filePath := "../../data/agent-instructions/booking-agent/quy-trinh-dat-lich.md"

	docs, _ := loadDocs(context.Background(), filePath)
	docsFormatted, err := convertToDocumentFormat(docs)
	assert.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		return len(docsFormatted) > 0
	}, "The len documents should be > 0")

	fmt.Println(docsFormatted[0])
}

func TestLoadBatchDocument(t *testing.T) {
	dataRoot := "../../data"
	filePaths, err := utils.GetListFiles(dataRoot)
	assert.NoError(t, err)

	docs, err := loadBatchDocs(context.Background(), filePaths)
	assert.NoError(t, err)
	assert.Condition(t, func() (success bool) {
		return len(docs) > 0
	}, "The len documents should be > 0")
}

func TestFormatBatchDocument(t *testing.T) {
	dataRoot := "../../data"
	filePaths, err := utils.GetListFiles(dataRoot)
	assert.NoError(t, err)

	docs, _ := loadBatchDocs(context.Background(), filePaths)
	docsFormatted, err := convertToDocumentFormat(docs)
	assert.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		return len(docsFormatted) > 0
	}, "The len documents should be > 0")

	fmt.Println(docsFormatted[0])
}
