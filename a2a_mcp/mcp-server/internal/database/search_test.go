package database

import (
	"context"
	ollama "learn-go/a2a_mcp/pkg/ollama"
	utils "learn-go/a2a_mcp/pkg/utils"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testTenantID = "11111111-1111-1111-1111-111111111111" // acme-corp from init-db.sql

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
	testDoc := &Document{
		TenantID:  testTenantID,
		Title:     "Test Document Without Embedding",
		Content:   "This document has no embedding vector and should not cause scan errors",
		Metadata:  map[string]interface{}{"test": true, "category": "integration-test"},
		Embedding: nil, // Explicitly no embedding
	}

	err := db.InsertDocument(ctx, testTenantID, testDoc)
	require.NoError(t, err, "Failed to insert test document")
	require.NotEmpty(t, testDoc.ID, "Document ID should be generated")

	// Now retrieve the document - this should NOT fail with NULL scan error
	retrieved, err := db.GetDocument(ctx, testTenantID, testDoc.ID)
	require.NoError(t, err, "Failed to retrieve document with NULL embedding")
	assert.NotNil(t, retrieved, "Retrieved document should not be nil")
	assert.Equal(t, testDoc.ID, retrieved.ID)
	assert.Equal(t, testDoc.Title, retrieved.Title)
	assert.Equal(t, testDoc.Content, retrieved.Content)
	assert.Nil(t, retrieved.Embedding, "Embedding should be nil for document without embedding")

	// Cleanup
	err = db.DeleteDocumentByID(ctx, testTenantID, testDoc.ID)
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
	testDoc := &Document{
		TenantID:  testTenantID,
		Title:     "Test Document With Embedding",
		Content:   "This document has an embedding vector",
		Metadata:  map[string]interface{}{"test": true, "category": "integration-test"},
		Embedding: embedding,
	}

	err := db.InsertDocument(ctx, testTenantID, testDoc)
	require.NoError(t, err, "Failed to insert test document")
	require.NotEmpty(t, testDoc.ID, "Document ID should be generated")

	// Retrieve the document
	retrieved, err := db.GetDocument(ctx, testTenantID, testDoc.ID)
	require.NoError(t, err, "Failed to retrieve document with embedding")
	assert.NotNil(t, retrieved, "Retrieved document should not be nil")
	assert.Equal(t, testDoc.ID, retrieved.ID)
	assert.Equal(t, testDoc.Title, retrieved.Title)
	assert.Equal(t, testDoc.Content, retrieved.Content)
	assert.NotNil(t, retrieved.Embedding, "Embedding should not be nil")
	assert.Equal(t, len(embedding), len(retrieved.Embedding), "Embedding dimension should match")

	// Cleanup
	err = db.DeleteDocumentByID(ctx, testTenantID, testDoc.ID)
	require.NoError(t, err, "Failed to delete test document")
}

func TestListDocuments_WithMixedEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create test documents with and without embeddings
	docs := []*Document{
		{
			TenantID:  testTenantID,
			Title:     "Doc 1 - No Embedding",
			Content:   "Content 1",
			Metadata:  map[string]interface{}{"test": true},
			Embedding: nil,
		},
		{
			TenantID:  testTenantID,
			Title:     "Doc 2 - Has Embedding",
			Content:   "Content 2",
			Metadata:  map[string]interface{}{"test": true},
			Embedding: make([]float32, 1024),
		},
		{
			TenantID:  testTenantID,
			Title:     "Doc 3 - No Embedding",
			Content:   "Content 3",
			Metadata:  map[string]interface{}{"test": true},
			Embedding: nil,
		},
	}

	// Insert all documents
	for _, doc := range docs {
		err := db.InsertDocument(ctx, testTenantID, doc)
		require.NoError(t, err, "Failed to insert document: "+doc.Title)
	}

	// List documents should handle mixed embeddings
	listed, err := db.ListDocuments(ctx, testTenantID, 10, 0)
	require.NoError(t, err, "Failed to list documents")
	assert.GreaterOrEqual(t, len(listed), 3, "Should have at least 3 documents")

	// Cleanup
	for _, doc := range docs {
		err = db.DeleteDocumentByID(ctx, testTenantID, doc.ID)
		require.NoError(t, err, "Failed to delete document: "+doc.Title)
	}
}

func TestSearchDocuments_WithNullEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Insert test document without embedding
	testDoc := &Document{
		TenantID:  testTenantID,
		Title:     "Security Policy Test",
		Content:   "Test content about security and authentication",
		Metadata:  map[string]interface{}{"test": true, "category": "security"},
		Embedding: nil,
	}

	err := db.InsertDocument(ctx, testTenantID, testDoc)
	require.NoError(t, err, "Failed to insert test document")

	// Search should work even with NULL embeddings
	results, err := db.SearchDocuments(ctx, testTenantID, "security", 10)
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
	err = db.DeleteDocumentByID(ctx, testTenantID, testDoc.ID)
	require.NoError(t, err, "Failed to delete test document")
}

func TestVectorSearch_SkipsNullEmbeddings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create query embedding
	queryEmbedding := make([]float32, 1536)
	for i := range queryEmbedding {
		queryEmbedding[i] = float32(i) * 0.001
	}

	// Vector search should only return documents with embeddings
	results, err := db.VectorSearch(ctx, testTenantID, queryEmbedding, 5)
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
	results, err := db.HybridSearch(ctx, testTenantID, params)
	require.NoError(t, err, "Hybrid search should not fail with NULL embeddings")
	assert.NotNil(t, results, "Results should not be nil")
	assert.Condition(t, func() (success bool) {
		return len(results) > 0
	}, "Hybrid search results should still appear")
}

func TestTenantIsolation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// 1. Verify RLS is enabled and user cannot bypass it
	tx, err := db.BeginTx(ctx, testTenantID)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// 2. Query the pg_class system table to ensure documents table has Row Level Security enabled.
	var rlsEnabled bool
	err = tx.QueryRow(ctx, "SELECT relrowsecurity FROM pg_class WHERE relname = 'documents'").Scan(&rlsEnabled) // 3. If the user has this permission (such as superuser privileges), RLS will be disabled.
	require.NoError(t, err, "Failed to check RLS status")
	require.True(t, rlsEnabled, "RLS must be enabled on documents table")

	// 3. If the user has this permission (such as superuser privileges), RLS will be disabled.
	var bypassRLS bool
	err = tx.QueryRow(ctx, "SELECT rolbypassrls FROM pg_roles WHERE rolname = current_user").Scan(&bypassRLS)
	require.NoError(t, err, "Failed to check user RLS bypass status")
	require.False(t, bypassRLS, "User must not have BYPASSRLS privilege for tenant isolation to work")

	tx.Commit(ctx)
	// 4. Insert sample data
	doc1 := &Document{
		TenantID: testTenantID,
		Title:    "Tenant 1 Document",
		Content:  "This belongs to tenant 1",
		Metadata: map[string]interface{}{"tenant": 1},
	}

	err = db.InsertDocument(ctx, testTenantID, doc1)
	require.NoError(t, err)

	// 5. Retrieve data from testTenantID
	retrieved, err := db.GetDocument(ctx, testTenantID, doc1.ID)
	require.NoError(t, err, "Should be able to retrieve document with correct tenant ID")
	assert.Equal(t, doc1.ID, retrieved.ID)

	// 6. Security Check (Attack/Unauthorized Access Test)
	otherTenantID := "22222222-2222-2222-2222-222222222222"
	_, err = db.GetDocument(ctx, otherTenantID, doc1.ID)
	assert.Error(t, err, "Should not be able to access document from different tenant")
	if err != nil {
		t.Logf("Expected error received: %v", err)
	}

	// Cleanup
	err = db.DeleteDocumentByID(ctx, testTenantID, doc1.ID)
	require.NoError(t, err)
}

func TestConcurrentRetrievals(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	var wg sync.WaitGroup

	// Get sample documents
	docs, err := db.ListDocuments(ctx, testTenantID, 20, 0)
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

			retrieved, err := db.GetDocument(ctx, testTenantID, docs[i].ID)
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
