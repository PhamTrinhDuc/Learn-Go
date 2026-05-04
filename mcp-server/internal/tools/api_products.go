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

type Product struct {
	ID       string  `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Category string  `json:"category" db:"category"`
	Price    float32 `json:"price_out" db:"price_out"`
	Stock    int     `json:"low_stock_threshold_retail" db:"low_stock_threshold_retail"`
}

type ProductTool struct {
	baseURL string
	client  *http.Client
}

func NewProductTool(baseURL string) *ProductTool {
	return &ProductTool{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Definition returns the tool definitions for Product API
func (t *ProductTool) ListProductDefinition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_products",
		Description: "List all products from chain store. Supports pagination.",
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

func (t *ProductTool) GetProductDefinition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_product_by_id",
		Description: "Get details of a specific product by their ID.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{
					"type":        "string",
					"description": "The unique ID of the product",
				},
			},
			"required": []string{"id"},
		},
	}
}

type ListProductsArgs struct {
	Page  float64 `json:"page"`
	Limit float64 `json:"limit"`
}

type GetProductArgs struct {
	ID string `json:"id"`
}

// Handlers
func (t *ProductTool) ListProductHandler(ctx context.Context, req *mcp.CallToolRequest, args ListProductsArgs) (*mcp.CallToolResult, any, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/products", t.baseURL))
	q := u.Query()

	if args.Page > 0 {
		q.Set("page", fmt.Sprintf("%.0f", args.Page))
	}
	if args.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%.0f", args.Limit))
	}
	u.RawQuery = q.Encode()
	apiURL := u.String()

	// 1. Tạo request mới với Context
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// 2. Inject Trace Context
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(httpReq.Header))

	// 3. Thực hiện gọi API
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
			Items []Product `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("failed to decode response: %v", err)}, nil
	}

	resultText := fmt.Sprintf("Found %d product(s):\n", len(apiResp.Data.Items))
	for _, s := range apiResp.Data.Items {
		resultText += fmt.Sprintf("- %s (Category: %s, Price: %.2f, Stock: %d)\n",
			s.Name, s.Category, s.Price, s.Stock)
	}

	return &mcp.CallToolResult{}, mcp.TextContent{Text: resultText}, nil
}

func (t *ProductTool) GetProductHandler(ctx context.Context, req *mcp.CallToolRequest, args GetProductArgs) (*mcp.CallToolResult, any, error) {
	apiURL := fmt.Sprintf("%s/products/%s", t.baseURL, args.ID)

	// 1. Tạo request mới với Context
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// 2. Inject Trace Context
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(httpReq.Header))

	// 3. Thực hiện gọi API
	resp, err := t.client.Do(httpReq)
	if err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("API request failed: %v", err)}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: "Product not found"}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("API returned error: %s", resp.Status)}, nil
	}

	var apiResp struct {
		Data Product `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &mcp.CallToolResult{IsError: true}, mcp.TextContent{Text: fmt.Sprintf("failed to decode response: %v", err)}, nil
	}

	s := apiResp.Data
	resultText := fmt.Sprintf("Product Details:\nID: %s\nName: %s\nCategory: %s\nPrice: %.2f\nStock: %d\n",
		s.ID, s.Name, s.Category, s.Price, s.Stock)

	return &mcp.CallToolResult{}, mcp.TextContent{Text: resultText}, nil
}
