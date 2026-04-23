package protocol

import (
	"encoding/json"
	"fmt"
)

const JSONRPCVersion = "2.0"

// Standard JSON-RPC error codes
const (
	InvalidRequest = -32600 // The JSON sent is not a valid Request object
	ParseError     = -32700 // Invalid JSON was received
)

// MCP-specific error codes (extending JSON-RPC)
const (
	AuthenticationRequired = -32001 // Authentication is required
	RateLimitExceeded      = -32004 // Rate limit exceeded
)

// Error represents a JSON-RPC 2.0 error object
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Request represents a JSON-RPC 2.0 request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"` // Can be string, number, or null
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

func NewRequest(id interface{}, code int, message string, data interface{}) (*Request, error) {
}

func NewResponse(id interface{}, result interface{}) *Response {

}

func NewErrorResponse(id interface{}, code int, message string, data interface{}) *Response {
	return &Response{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

func (r *Request) Validate() error {
	if r.JSONRPC != JSONRPCVersion {
		return fmt.Errorf("Invalid jsonrpc version: expected %s, got %s", JSONRPCVersion, r.JSONRPC)
	}

	if r.Method == "" {
		return fmt.Errorf("Method is required")
	}
	// ID can be string, number, or null - we accept all in interface{}
	return nil
}
