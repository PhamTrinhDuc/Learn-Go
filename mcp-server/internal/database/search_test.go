package database

import (
	"context"
	ollama "mcp-server/internal/llm"
	utils "mcp-server/internal/utils"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDBConfig() DBConfig {
	return DBConfig{
		Host:     utils.GetEnvString("DB_HOST", "localhost"),
		Port:     utils.GetEnvInt("DB_PORT", 5433),
		User:     utils.GetEnvString("DB_USER", "app_user"), // Use app_user for RLS enforcement
		Password: utils.GetEnvString("DB_PASSWORD", "mcp_password"),
		DBName:   utils.GetEnvString("DB_NAME", "mcp_db"),
		SSLMode:  utils.GetEnvString("DB_SSLMODE", "disable"),
		MaxConns: int32(utils.GetEnvInt("MAX_CONNS", 10)),
		MinConns: int32(utils.GetEnvInt("MAX_CONNS", 2)),
	}
}

func getTestModelConfig() ollama.Config {
	return ollama.Config{
		BaseURL:    utils.GetEnvString("OLLAMA_URL", "http://localhost:11434"),
		LLMModel:   utils.GetEnvString("LLM_MODEL", "qwen2.5:0.5b"),
		EmbedModel: utils.GetEnvString("EMBED_MODEL", "qwen3-embedding:0.6b"),
	}
}

func setupTestDB(t *testing.T) *DB {
	cfg := getTestDBConfig()
	db, err := NewDB(context.Background(), cfg)
	require.NoError(t, err, "Failed to connect to test database")

	// model, err := ollama.NewClient(getTestModelConfig())
	// require.NoError(t, err, "failed to to setup Ollama model")
	return db
}

func TestGetDocument_WithNullEmbedding(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Insert a test document WITHOUT embedding
	testDoc := &KnowledgeBase{
		//
		Title:     "Test Document Without Embedding",
		Content:   "This document has no embedding vector and should not cause scan errors",
		Metadata:  map[string]interface{}{"test": true, "category": "integration-test"},
		Embedding: nil, // Explicitly no embedding
	}

	err := db.InsertDocument(ctx, testDoc)
	require.NoError(t, err, "Failed to insert test document")
	require.NotEmpty(t, testDoc.ID, "Document ID should be generated")

	// Now retrieve the document - this should NOT fail with NULL scan error
	retrieved, err := db.GetDocument(ctx, testDoc.ID)
	require.NoError(t, err, "Failed to retrieve document with NULL embedding")
	assert.NotNil(t, retrieved, "Retrieved document should not be nil")
	assert.Equal(t, testDoc.ID, retrieved.ID)
	assert.Equal(t, testDoc.Title, retrieved.Title)
	assert.Equal(t, testDoc.Content, retrieved.Content)
	assert.Nil(t, retrieved.Embedding, "Embedding should be nil for document without embedding")

	// Cleanup
	err = db.DeleteDocumentByID(ctx, testDoc.ID)
	require.NoError(t, err, "Failed to delete test document")
}

func TestGetDocument_WithEmbedding(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a test embedding vector (1536 dimensions for OpenAI ada-002)
	embedding := make([]float32, 1024)
	for i := range embedding {
		embedding[i] = float32(i) * 0.001
	}

	// Insert a test document WITH embedding
	testDoc := &KnowledgeBase{
		//
		Title:     "Test Document With Embedding",
		Content:   "This document has an embedding vector",
		Metadata:  map[string]interface{}{"test": true, "category": "integration-test"},
		Embedding: embedding,
	}

	err := db.InsertDocument(ctx, testDoc)
	require.NoError(t, err, "Failed to insert test document")
	require.NotEmpty(t, testDoc.ID, "Document ID should be generated")

	// Retrieve the document
	retrieved, err := db.GetDocument(ctx, testDoc.ID)
	require.NoError(t, err, "Failed to retrieve document with embedding")
	assert.NotNil(t, retrieved, "Retrieved document should not be nil")
	assert.Equal(t, testDoc.ID, retrieved.ID)
	assert.Equal(t, testDoc.Title, retrieved.Title)
	assert.Equal(t, testDoc.Content, retrieved.Content)
	assert.NotNil(t, retrieved.Embedding, "Embedding should not be nil")
	assert.Equal(t, len(embedding), len(retrieved.Embedding), "Embedding dimension should match")

	// Cleanup
	err = db.DeleteDocumentByID(ctx, testDoc.ID)
	require.NoError(t, err, "Failed to delete test document")
}

func TestListDocuments_WithMixedEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create test documents with and without embeddings
	docs := []*KnowledgeBase{
		{
			//
			Title:     "Doc 1 - No Embedding",
			Content:   "Content 1",
			Metadata:  map[string]interface{}{"test": true},
			Embedding: nil,
		},
		{
			//
			Title:     "Doc 2 - Has Embedding",
			Content:   "Content 2",
			Metadata:  map[string]interface{}{"test": true},
			Embedding: make([]float32, 1024),
		},
		{
			//
			Title:     "Doc 3 - No Embedding",
			Content:   "Content 3",
			Metadata:  map[string]interface{}{"test": true},
			Embedding: nil,
		},
	}

	// Insert all documents
	for _, doc := range docs {
		err := db.InsertDocument(ctx, doc)
		require.NoError(t, err, "Failed to insert document: "+doc.Title)
	}

	// List documents should handle mixed embeddings
	listed, err := db.ListDocuments(ctx, 10, 0)
	require.NoError(t, err, "Failed to list documents")
	assert.GreaterOrEqual(t, len(listed), 3, "Should have at least 3 documents")

	// Cleanup
	for _, doc := range docs {
		err = db.DeleteDocumentByID(ctx, doc.ID)
		require.NoError(t, err, "Failed to delete document: "+doc.Title)
	}
}

func TestSearchDocuments_WithNullEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Insert test document without embedding
	testDoc := &KnowledgeBase{
		//
		Title:     "Security Policy Test",
		Content:   "Test content about security and authentication",
		Metadata:  map[string]interface{}{"test": true, "category": "security"},
		Embedding: nil,
	}

	err := db.InsertDocument(ctx, testDoc)
	require.NoError(t, err, "Failed to insert test document")

	// Search should work even with NULL embeddings
	results, err := db.SearchDocuments(ctx, "security", 10)
	require.NoError(t, err, "Failed to search documents")
	assert.GreaterOrEqual(t, len(results), 1, "Should find at least one document")

	// Verify our test document is in results
	found := false
	for _, doc := range results {
		if doc.ID == testDoc.ID {
			found = true
			assert.Equal(t, testDoc.Title, doc.Title)
			break
		}
	}
	assert.True(t, found, "Should find our test document in search results")

	// Cleanup
	err = db.DeleteDocumentByID(ctx, testDoc.ID)
	require.NoError(t, err, "Failed to delete test document")
}

func TestVectorSearch_SkipsNullEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create query embedding
	queryEmbedding := make([]float32, 1535)
	for i := range queryEmbedding {
		queryEmbedding[i] = float32(i) * 0.001
	}

	// Vector search should only return documents with embeddings
	results, err := db.VectorSearch(ctx, queryEmbedding, 5)
	require.NoError(t, err, "Vector search should not fail")

	// All returned documents should have embeddings
	for _, result := range results {
		assert.NotNil(t, result.Document.Embedding, "Vector search should only return docs with embeddings")
	}
}

func TestHybridSearch_HandlesNullEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create query embedding
	queryEmbedding := make([]float32, 1024)
	for i := range queryEmbedding {
		queryEmbedding[i] = 0.1
	}

	params := HybridSearchParams{
		Query:        "security policy",
		Embedding:    queryEmbedding,
		Limit:        10,
		BM25Weight:   0.5,
		VectorWeight: 0.5,
		MinBM25Score: 0.0,
		MinVectorSim: 0.0,
	}

	// Hybrid search should handle documents without embeddings gracefully
	results, err := db.HybridSearch(ctx, params)
	require.NoError(t, err, "Hybrid search should not fail with NULL embeddings")
	assert.NotNil(t, results, "Results should not be nil")
	assert.Condition(t, func() (success bool) {
		return len(results) > 0
	}, "Hybrid search results should still appear")
}

func TestConcurrentRetrievals(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	var wg sync.WaitGroup

	// Get sample documents
	docs, err := db.ListDocuments(ctx, 20, 0)
	require.NoError(t, err)
	require.NotEmpty(t, docs)

	// Perform concurrent retrievals
	numWorkers := 10
	sem := make(chan struct{}, numWorkers)

	for i := range len(docs) {
		wg.Add(1) // Đếm thêm 1 việc cần làm
		// Đợi lấy "vé" từ semaphore, nếu đủ 10 người rồi thì sẽ kẹt lại ở đây
		sem <- struct{}{}

		go func(workerID int) {
			defer wg.Done() // Đảm bảo dù lỗi hay không cũng sẽ trừ máy đếm khi xong

			defer func() {
				<-sem
			}() // Xong việc thì trả lại vé cho người sau. Không để ngoài func vì phải chờ for xong mới thực thi => không ai trả vé

			retrieved, err := db.GetDocument(ctx, docs[i].ID)
			// t.Logf("received ID %s for document", retrieved.ID)
			if err != nil {
				t.Errorf("Worker %d failed to retrieve document: %v", workerID, err)
				return
			}
			if retrieved == nil {
				t.Errorf("Worker %d got nil document", workerID)
				return
			}
		}(i)
	}
	// 4. Đợi tất cả hoàn thành
	wg.Wait()
	t.Log("✓ All concurrent retrievals completed successfully")
}
