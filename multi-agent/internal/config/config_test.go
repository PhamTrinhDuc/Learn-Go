package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAgentConfig(t *testing.T) {
	// Load config
	configFile := "../../config.yaml"
	bookingAgent, err := LoadAgentConfig(configFile, "booking_agent")
	assert.NoError(t, err)
	qaAgent, err := LoadAgentConfig(configFile, "qa_agent")
	assert.NoError(t, err)

	// Kiểm tra booking_agent
	assert.Equal(t, "booking_agent", bookingAgent.Name)
	assert.Contains(t, bookingAgent.Tools, "get_available_slots")
	assert.Contains(t, bookingAgent.Tools, "booking")

	// Kiểm tra qa_agent
	assert.Equal(t, "qa_agent", qaAgent.Name)
	assert.Contains(t, qaAgent.Tools, "search_documents")
	assert.Contains(t, qaAgent.Tools, "get_products")
}

func TestLoadAppConfig(t *testing.T) {
	configFile := "../../config.yaml"
	config, err := LoadAppConfig(configFile)
	assert.NoError(t, err)
	assert.Equal(t, "gemini-flash-2.5", config.Models["gemini"])
	assert.Equal(t, "openai-4o-mini", config.Models["openai"])
	assert.Equal(t, "http://localhost:8081", config.McpServer)
}
