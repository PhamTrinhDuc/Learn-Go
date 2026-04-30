package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

func NewMCPTool(mcpServerURL string, allowedTools []string, approvalTools []string) (tool.Toolset, error) {
	// 1. Thiết lập Transport kết nối với MCP Server (Sử dụng SSE transport qua HTTP)
	// Khởi tạo SSE client để giao tiếp với mcp-server
	transport := &mcp.SSEClientTransport{Endpoint: mcpServerURL}

	// 2. Khởi tạo Toolset từ MCP
	// mcptoolset sẽ fetch toàn bộ các tool từ mcp-server
	mcpTools, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
		// Chỉ cho phép các tools được định nghĩa trong allowedTools
		ToolFilter: tool.StringPredicate(allowedTools),
		// Human-in-the-loop: Yêu cầu xác nhận (approve) cho một số action nhạy cảm
		RequireConfirmationProvider: func(toolName string, toolInput any) bool {
			for _, tool := range approvalTools {
				if tool == toolName {
					return true
				}
			}
			return false
		},
	})
	if err != nil {
		return nil, err
	}
	return mcpTools, nil
}
