package index

import (
	"context"
	"fmt"

	database "mcp-server/internal/database"
	ollama "mcp-server/internal/llm"
)

func ingestToVDB(ctx context.Context, model *ollama.Client, db *database.DB, filePath string, tenantID string) error {
	docs, err := loadDocs(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to load documents: %w", err)
	}

	docFormatted, err := convertToDocumentFormat(docs, tenantID)
	if err != nil {
		return fmt.Errorf("failed to convert document: %w", err)
	}

	// 3. AI sinh Vector cho toàn bộ chunks
	// Trích xuất text từ danh sách docs để gửi cho AI
	texts := make([]string, len(docFormatted))
	for i, d := range docFormatted {
		texts[i] = d.Content
	}

	fmt.Printf("[*] Đang sinh Embedding cho %d đoạn văn bản theo từng đợt...\n", len(texts))
	const batchSize = 50
	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batchVectors, err := model.GenerateEmbeddings(ctx, texts[i:end])
		if err != nil {
			return fmt.Errorf("failed to generate embeddings (batch %d-%d): %w", i, end, err)
		}

		// Gán vector vào từng document
		for j, v := range batchVectors {
			docFormatted[i+j].Embedding = v
		}
		fmt.Printf("    - Đã xong: %d/%d\n", end, len(texts))
	}

	// Dùng InsertBatchingDocument vì docFormatted là một slice (danh sách)
	err = db.InsertBatchingDocument(ctx, tenantID, docFormatted, 6)
	if err != nil {
		return fmt.Errorf("failed to insert docs: %w", err)
	}
	return nil
}

func ingestBatchingToVDB(ctx context.Context, model *ollama.Client, db *database.DB, filePaths []string, tenantID string) error {
	// 1. Load và Split song song (Tận dụng goroutines đã viết ở spliiter.go)
	docs, err := loadBatchDocs(ctx, filePaths)
	if err != nil {
		return fmt.Errorf("failed to load batch docs: %w", err)
	}

	// 2. Format dữ liệu
	docFormatted, err := convertToDocumentFormat(docs, tenantID)
	if err != nil {
		return fmt.Errorf("failed to convert documents: %w", err)
	}

	// 3. AI sinh Vector cho toàn bộ chunks
	// Trích xuất text từ danh sách docs để gửi cho AI
	texts := make([]string, len(docFormatted))
	for i, d := range docFormatted {
		texts[i] = d.Content
	}

	fmt.Printf("[*] Đang sinh Embedding cho %d đoạn văn bản theo từng đợt...\n", len(texts))
	const batchSize = 50
	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batchVectors, err := model.GenerateEmbeddings(ctx, texts[i:end])
		if err != nil {
			return fmt.Errorf("failed to generate embeddings (batch %d-%d): %w", i, end, err)
		}

		// Gán vector vào từng document
		for j, v := range batchVectors {
			docFormatted[i+j].Embedding = v
		}
		fmt.Printf("    - Đã xong: %d/%d\n", end, len(texts))
	}

	// 4. Insert một cú duy nhất vào DB (Batching)
	err = db.InsertBatchingDocument(ctx, tenantID, docFormatted, 6)
	if err != nil {
		return fmt.Errorf("failed to insert batch to database: %w", err)
	}

	return nil
}
