package agents

import (
	"context"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	config "multi-agent/config"
)

// NewBookingAgent tạo template cho Booking Agent, kết nối tới mcp-server.
func NewSubAgent(ctx context.Context, cfg *config.AgentConfig, model model.LLM, mcpTools tool.Toolset) (agent.Agent, error) {

	// 1. Tạo LLM Agent với role và tools
	agent, err := llmagent.New(llmagent.Config{
		Name:        cfg.Name,
		Model:       model,
		Description: cfg.Description,
		Instruction: cfg.Instruction,
		Toolsets:    []tool.Toolset{mcpTools},
	})
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func NewOrscheratorAgent(ctx context.Context, cfg *config.AgentConfig, model model.LLM, subAgents []agent.Agent) (agent.Agent, error) {
	// 1. Tạo LLM Agent với role và tools
	agent, err := llmagent.New(llmagent.Config{
		Name:        cfg.Name,
		Model:       model,
		Description: cfg.Description,
		Instruction: cfg.Instruction,
		SubAgents:   subAgents,
	})
	if err != nil {
		return nil, err
	}

	return agent, nil
}
