package tools

import (
	"context"
	"fmt"
	"mcp-server/internal/auth"
	"mcp-server/internal/database"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchTool struct {
	db database.Store
}

func NewSearchTool(db database.Store) *SearchTool {
	return &SearchTool{db: db}
}

func (t *SearchTool) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "search_documents",
		Description: "Search documents by text query.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "the search query text",
				},
				"limit": map[string]any{
					"type":        "number",
					"description": "Maximum number of results to return",
				},
			},
			"required": []string{"query"},
		},
	}
}

type SearchArgs struct {
	Query string  `json:"query"`
	Limit float64 `json:"limit"`
}

func (t *SearchTool) Handler(ctx context.Context, req *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, any, error) {
	tenantID, err := auth.ExtractTenantID(ctx)
	if err != nil {
		res := &mcp.CallToolResult{IsError: true}
		return res, mcp.TextContent{Text: fmt.Sprintf("authentication required: %v", err)}, nil
	}

	limit := int(args.Limit)
	if limit <= 0 {
		limit = 10
	}

	documents, err := t.db.SearchDocuments(ctx, tenantID, args.Query, limit)
	if err != nil {
		res := &mcp.CallToolResult{IsError: true}
		return res, mcp.TextContent{Text: fmt.Sprintf("search failed: %v", err)}, nil
	}

	var resultText string
	if len(documents) == 0 {
		resultText = fmt.Sprintf("No documents found matching query: %s", args.Query)
	} else {
		resultText = fmt.Sprintf("Found %d document(s)\n", len(documents))
		for i, doc := range documents {
			resultText += fmt.Sprintf("%d. %s\n", i+1, doc.Title)
		}
	}

	return &mcp.CallToolResult{}, mcp.TextContent{Text: resultText}, nil
}
