package main

// https://github.com/achetronic/adk-utils-go.git
import (
	"context"
	"fmt"
	"log"

	"multi-agent/internal/agents"
	config "multi-agent/internal/config"
	mymcp "multi-agent/internal/mcp"
	"multi-agent/internal/provider"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

const (
	appName = "salon_chain"
	userID  = "demo_user"
)

func runAgent(ctx context.Context, runnr *runner.Runner, sessionID string, input string) string {
	userMsg := genai.NewContentFromText(input, genai.RoleUser)

	var responseText string
	for event, err := range runnr.Run(ctx, userID, sessionID, userMsg, agent.RunConfig{}) {
		if err != nil {
			log.Printf("Error: %v", err)
			break
		}
		if event.ErrorCode != "" {
			log.Printf("Event error: %s - %s", event.ErrorCode, event.ErrorMessage)
			break
		}
		if event.Content != nil && len(event.Content.Parts) > 0 {
			responseText += event.Content.Parts[0].Text
		}
	}

	return responseText
}

func main() {
	ctx := context.Background()
	cfgPath := "../config.yaml"

	// 1. Load App Config (contains all agent definitions)
	appCfg, err := config.LoadAppConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load app config: %v", err)
	}

	// 2. Init Shared Resources
	// Shared LLM model
	llm, err := provider.NewGeminiLLM(ctx, "gemini-2.0-flash") // Sử dụng model mặc định hoặc từ config
	if err != nil {
		log.Fatalf("Failed to create LLM model: %v", err)
	}

	// Shared MCP Transport
	transport := &mcp.SSEClientTransport{Endpoint: appCfg.McpServer}

	// 3. Initialize Agents in Parallel
	registry := agents.NewRegistry()
	g, ctx := errgroup.WithContext(ctx)

	for agentName, agentCfg := range appCfg.Agents {
		if agentName == "orschestrator_agent" {
			continue
		}
		// Capture loop variables
		name := agentName
		cfg := agentCfg

		g.Go(func() error {
			// Initialize MCP tools for this agent
			mcpToolset, err := mymcp.NewMCPTool(transport, cfg.AllowedTools, cfg.ApprovedTools)
			if err != nil {
				return fmt.Errorf("failed to create MCP tools for %s: %w", name, err)
			}

			// Create the Agent instance
			newAgent, err := agents.NewSubAgent(ctx, &cfg, llm, mcpToolset)
			if err != nil {
				return fmt.Errorf("failed to create agent %s: %w", name, err)
			}

			// Register in our global registry
			registry.Register(name, newAgent)
			log.Printf("Successfully initialized sub-agent: %s", name)
			return nil
		})
	}

	// Wait for all sub-agents to be ready
	if err := g.Wait(); err != nil {
		log.Fatalf("Sub-agent initialization failed: %v", err)
	}

	// 4. Init orschestrator agent
	var targetAgent agent.Agent
	if orschestratorCfg, ok := appCfg.Agents["orschestrator_agent"]; ok {
		var err error
		targetAgent, err = agents.NewOrscheratorAgent(ctx, &orschestratorCfg, llm, registry.GetAgents())
		if err != nil {
			log.Fatalf("Failed to create orschestrator agent: %v", err)
		}
		log.Printf("Orschestrator agent initialized: %s", orschestratorCfg.Name)
	} else {
		// Fallback to any agent if orschestrator_agent is not found
		log.Printf("Orschestrator agent not found, falling back to first available agent")
		names := registry.ListNames()
		if len(names) > 0 {
			targetAgent, _ = registry.Get(names[0])
		} else {
			log.Fatal("No agents registered to launch")
		}
	}

	// 5. Launch the application
	runnr, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          targetAgent,
		SessionService: session.InMemoryService(),
		MemoryService:  memory.InMemoryService(), // Enables automatic memory persistence
	})

	userInput := "Hello"
	response := runAgent(ctx, runnr, uuid.NewString(), userInput)
	fmt.Println(response)

}
