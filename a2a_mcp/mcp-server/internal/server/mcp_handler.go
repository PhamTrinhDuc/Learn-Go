package server

import (
	"mcp-server/internal/database"
	"mcp-server/internal/tools"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewSSEHandler(db database.Store) http.Handler {
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "a2a-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	searchTool := tools.NewSearchTool(db)
	hybridTool := tools.NewHybridSearchTool(db)

	mcp.AddTool(s, searchTool.Definition(), searchTool.Handler)
	mcp.AddTool(s, hybridTool.Definition(), hybridTool.Handler)

	// Create and return the SSE HTTP Handler using the discovered function name
	return mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		return s
	}, nil)
}
