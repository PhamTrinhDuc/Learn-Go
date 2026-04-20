package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"learn-go/a2a_mcp/mcp-server/internal/auth"
	"learn-go/a2a_mcp/mcp-server/internal/database"
	"learn-go/a2a_mcp/mcp-server/internal/protocol"
	"learn-go/a2a_mcp/pkg/ollama"
)

type SearchTool struct {
	db    database.Store
	model *ollama.Client
}

type SearchParams struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

func NewSearchTool(db database.Store, model *ollama.Client) *SearchTool {
	return &SearchTool{db: db, model: model}
}

// Definition returns the tool definition for MCP
func (t *SearchTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        "search documents",
		Description: "Search documents by text query. Searches across title, content, and metadata fields.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "the search query text",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results to return (defaults: 10, max: 50)",
				},
			},
			"required": []string{"query"},
		},
	}
}

func (t *SearchTool) Execute(
	ctx context.Context,
	args map[string]interface{}) (protocol.ToolCallResult, error) {
	// 1. Extract tenant ID from context
	tenantID, err := auth.ExtractTenantID(ctx)
	if err != nil {
		return protocol.ToolCallResult{IsError: true},
			fmt.Errorf("authentication required: %w", err)
	}

	// 2. Parse json parameters
	argsJson, err := json.Marshal(args)
	if err != nil {
		return protocol.ToolCallResult{IsError: true},
			fmt.Errorf("failed to parse json arguments")
	}

	// 3. Valid params
	var params SearchParams
	if err := json.Unmarshal(argsJson, &params); err != nil {
		return protocol.ToolCallResult{IsError: true},
			fmt.Errorf("failed to convert args to schema: %w", err)
	}

	if params.Query == "" {
		return protocol.ToolCallResult{IsError: true},
			fmt.Errorf("query is required")
	}

	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 50
	}

	// 4. Generate embedding for the query
	embeddings, err := t.model.GenerateEmbeddings(ctx, []string{params.Query})
	if err != nil {
		return protocol.ToolCallResult{IsError: true},
			fmt.Errorf("failed to generate embedding for query: %w", err)
	}

	// 5. Execute Hybrid Search
	searchParams := database.HybridSearchParams{
		Query:        params.Query,
		Embedding:    embeddings[0],
		Limit:        params.Limit,
		BM25Weight:   0.7, // Thiên hướng tìm từ khóa chính xác
		VectorWeight: 0.3, // Kết hợp với ý nghĩa vector
	}
	results, err := t.db.HybridSearch(ctx, tenantID, searchParams)
	if err != nil {
		return protocol.ToolCallResult{IsError: true},
			fmt.Errorf("failed to perform hybrid search: %w", err)
	}

	// Formats results for LLM
	var resultText string
	if len(results) == 0 {
		resultText = fmt.Sprintf("No documents found matching query: %s", params.Query)
	} else {
		resultText = fmt.Sprintf("Found %d document(s) matching query: %s\n\n", len(results), params.Query)
		for i, res := range results {
			doc := res.Document
			resultText += fmt.Sprintf("Document %d (Score: %.4f):\n", i+1, res.CombinedScore)
			resultText += fmt.Sprintf("  ID: %s\n", doc.ID)
			resultText += fmt.Sprintf("  Title: %s\n", doc.Title)
			resultText += fmt.Sprintf("  Content Preview: %.200s...\n", doc.Content)
			if doc.Metadata != nil {
				metadataJSON, _ := json.Marshal(doc.Metadata)
				resultText += fmt.Sprintf("  Metadata: %s\n", string(metadataJSON))
			}
			resultText += "\n"
		}
	}

	return protocol.ToolCallResult{
		Content: []protocol.ContentBlock{
			{
				Type: "text",
				Text: resultText,
			},
		},
		IsError: false,
	}, nil
}
