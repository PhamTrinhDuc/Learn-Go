package server

import (
	"encoding/json"
	"io"
	"learn-go/a2a_mcp/mcp-server/internal/observability"
	"learn-go/a2a_mcp/mcp-server/internal/protocol"
	"learn-go/a2a_mcp/mcp-server/internal/tools"
	"net/http"
)

type MCPHandler struct {
	toolRegistry *tools.Registry
	telemetry    *observability.Telemetry
}

// NewMCPHandler creates a new MCP handler
func NewMCPHandler(toolRegistry *tools.Registry, telemetry *observability.Telemetry) *MCPHandler {
	return &MCPHandler{
		toolRegistry: toolRegistry,
		telemetry:    telemetry,
	}
}

// ServeHTTP implements http.Handler
func (h *MCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ctx := context.Background()
	// startTime := time.Now()

	// 1. accept only POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendErrorResponse(w, nil, protocol.ParseError, "Failed to read request body")
	}

	// 3. Parse JSON-RPC Request
	var req protocol.Request
	if err := json.Unmarshal(body, &req); err != nil {
		h.sendErrorResponse(w, nil, protocol.ParseError, "Invald request JSON")
		return
	}

	// 4. Validate request
	if err := req.Validate(); err != nil {
		h.sendErrorResponse(w, nil, protocol.InvalidRequest, err.Error())
		return
	}

	// 5. Start tracing span
	

}

func (h *MCPHandler) sendResponse(w http.ResponseWriter, response *protocol.Response) {
	w.Header().Set("Content-Type", "application/json")

	if response.Error != nil {
		switch response.Error.Code {
		case protocol.ParseError:
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *MCPHandler) sendErrorResponse(w http.ResponseWriter, id interface{}, code int, message string) {
	response := protocol.NewErrorResponse(id, code, message, nil)
	h.sendResponse(w, response)
}
