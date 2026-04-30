package main

import (
	"context"
	"fmt"
	"log"
	"multi-agent/internal/agents"
	config "multi-agent/internal/config"
	"multi-agent/internal/mcp"
	"multi-agent/internal/provider"
	"os"

	"google.golang.org/adk/agent"

	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
)

func main() {
	cfgPath := "../config.yaml"
	ctx := context.Background()
	cfg, err := config.LoadAppConfig(cfgPath)

	cfgAgent, err := config.LoadAgentConfig(cfgPath, "qa_agent")

	if err != nil {
		fmt.Println("Failed to init config app: ", err)
		return
	}
	mcpTools, err := mcp.NewMCPTool(cfg.McpServer, cfgAgent.Tools, []string{})
	if err != nil {
		fmt.Println("Failed to create MCP Tools")
		return
	}

	llm, err := provider.NewGeminiLLM(ctx, "gemini-flash-2.5")
	if err != nil {
		fmt.Println("Failed to create LLM model: ", err)
		return
	}

	myAgent, err := agents.NewAgent(ctx, cfgAgent, llm, mcpTools)
	if err != nil {
		fmt.Println("Failed to create new Agent: ", err)
		return
	}

	launchCfg := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(myAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, launchCfg, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
