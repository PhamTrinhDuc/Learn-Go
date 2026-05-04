package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Branch struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Address      string `json:"address" db:"address"`
	Phone        string `json:"phone" db:"phone"`
	OpeningHours string `json:"opening_hours" db:"opening_hours"`
	IsActive     bool   `json:"is_active" db:"is_active"`
}

type BranchTool struct {
	baseURL string
	client  *http.Client
}

func NewBranchTool(baseURL string) *BranchTool {
	return &BranchTool{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Definition returns the tool definitions for Stylist API
func (t *BranchTool) ListBranchDefinition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_branches",
		Description: "List all branchs from chain store. Supports pagination.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"page": map[string]any{
					"type":        "number",
					"description": "Page number (default 1)",
				},
				"limit": map[string]any{
					"type":        "number",
					"description": "Number of items per page (default 10)",
				},
			},
		},
	}
}

func (t *BranchTool) GetBranchDefinition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_branch_by_id",
		Description: "Get details of a specific branch by their ID.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{
					"type":        "string",
					"description": "The unique ID of the branch",
				},
			},
			"required": []string{"id"},
		},
	}
}

type ListBranchsArgs struct {
	Page  float64 `json:"page"`
	Limit float64 `json:"limit"`
}

type GetBranchArgs struct {
	ID string `json:"id"`
}

// Handlers
func (t *BranchTool) ListBranchHandler(ctx context.Context, req *mcp.CallToolRequest, args ListBranchsArgs) (*mcp.CallToolResult, any, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/branches", t.baseURL))
	q := u.Query()

	if args.Page > 0 {
		q.Set("page", fmt.Sprintf("%.0f", args.Page))
	}
	if args.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%.0f", args.Limit))
	}

	u.RawQuery = q.Encode()
	apiURL := u.String()

	// 1. Tạo request mới với Context để mang theo span hiện tại
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// 2. Inject Trace Context vào Header
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(httpReq.Header))

	// 3. Thực hiện gọi API bằng request đã được inject header
	resp, err := t.client.Do(httpReq)
	if err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("API request failed: %v", err)}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("API returned error: %s", resp.Status)}, nil
	}

	var apiResp struct {
		Data struct {
			Items []Branch `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("failed to decode response: %v", err)}, nil
	}

	resultText := fmt.Sprintf("Found %d branch(s):\n", len(apiResp.Data.Items))
	for _, s := range apiResp.Data.Items {
		resultText += fmt.Sprintf("Branch Details:\nName: %s\nAddress: %s\nPhone: %s\nOpenhours: %s\nActive: %v\n",
			s.Name, s.Address, s.Phone, s.OpeningHours, s.IsActive)
	}

	return &mcp.CallToolResult{}, mcp.TextContent{Text: resultText}, nil
}

func (t *BranchTool) GetBranchHandler(ctx context.Context, req *mcp.CallToolRequest, args GetBranchArgs) (*mcp.CallToolResult, any, error) {
	apiURL := fmt.Sprintf("%s/branches/%s", t.baseURL, args.ID)

	// 1. Tạo request mới với Context
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// 2. Inject Trace Context vào Header
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(httpReq.Header))

	// 3. Thực hiện gọi API
	resp, err := t.client.Do(httpReq)
	if err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("API request failed: %v", err)}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: "Branch not found"}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("API returned error: %s", resp.Status)}, nil
	}

	var apiResp struct {
		Data Branch `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("failed to decode response: %v", err)}, nil
	}

	s := apiResp.Data
	resultText := fmt.Sprintf("Branch Details:\nName: %s\nAddress: %s\nPhone: %s\nOpenhours: %s\nActive: %v\n",
		s.Name, s.Address, s.Phone, s.OpeningHours, s.IsActive)

	return &mcp.CallToolResult{}, mcp.TextContent{Text: resultText}, nil
}
