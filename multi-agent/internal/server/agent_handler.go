package main

// https://github.com/achetronic/adk-utils-go.git
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"multi-agent/internal/agents"
	config "multi-agent/internal/config"
	mymcp "multi-agent/internal/mcp"
	"multi-agent/internal/observability"
	"multi-agent/internal/provider/gemini"
	"multi-agent/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

// headerTransport là một RoundTripper tùy chỉnh để chèn thêm header vào mọi request
type headerTransport struct {
	base   http.RoundTripper
	header string
	value  string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone request để không làm ảnh hưởng đến bản gốc
	newReq := req.Clone(req.Context())
	newReq.Header.Set(t.header, t.value)

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(req.Context(), propagation.HeaderCarrier(newReq.Header))
	return t.base.RoundTrip(newReq)
}

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

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type ChatResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type AgentServer struct {
	Runner         *runner.Runner
	SessionService session.Service
	Config         *config.AppConfig
	Telemetry      *observability.Telemetry
}

func loadConfig() observability.Config {
	return observability.Config{
		ServiceName:    utils.GetEnvString("OTEL_SERVICE", "agent-server"),
		ServiceVersion: utils.GetEnvString("OTEL_VERSION", "1.0.0"),
		Environment:    utils.GetEnvString("ENVIRONMENT", "development"),
		OTLPEndpoint:   utils.GetEnvString("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
		SamplingRate:   utils.GetEnvFloat("OTEL_TRACES_SAMPLER_ARG", 1.0),
		EnableTracing:  utils.GetEnvBool("OTEL_ENABLE_TRACING", true),
		EnableMetrics:  utils.GetEnvBool("OTEL_ENABLE_METRICS", true),
	}
}

func NewAgentServer(ctx context.Context, cfgPath string) (*AgentServer, error) {
	// 1. Load App Config (contains all agent definitions)
	appCfg, err := config.LoadAppConfig(cfgPath)
	if err != nil {
		log.Printf("Failed to load app config: %v", err)
	}

	// 2. Init Shared Resources
	// Shared LLM model
	llm, err := gemini.NewGeminiLLM(ctx, "gemini-2.5-flash") // Gemini 2.0 Flash is stable
	if err != nil {
		log.Printf("Failed to create LLM model: %v", err)
	}

	mcpToken := utils.GetEnvString("AUTH_TOKEN", "")

	// Shared MCP Transport
	transport := &mcp.SSEClientTransport{
		Endpoint: appCfg.McpServer,
		HTTPClient: &http.Client{
			Transport: &headerTransport{
				base:   http.DefaultTransport,
				header: "Authorization",
				value:  "Bearer " + mcpToken,
			},
		},
	}

	// 3. Initialize Agents in Parallel
	registry := agents.NewRegistry()
	g, groupCtx := errgroup.WithContext(ctx)

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
			// Use groupCtx for initialization
			newAgent, err := agents.NewSubAgent(groupCtx, &cfg, llm, mcpToolset)
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
		log.Printf("Sub-agent initialization failed: %v", err)
		return nil, fmt.Errorf("failed to init sub-agents: %w", err)
	}

	// 4. Init orschestrator agent
	var targetAgent agent.Agent
	if orschestratorCfg, ok := appCfg.Agents["orschestrator_agent"]; ok {
		var err error
		targetAgent, err = agents.NewOrscheratorAgent(ctx, &orschestratorCfg, llm, registry.GetAgents())
		if err != nil {
			log.Printf("Failed to create orschestrator agent: %v", err)
			return nil, fmt.Errorf("Failed to create orschestrator agent: %w", err)
		}
		log.Printf("Orschestrator agent initialized: %s", orschestratorCfg.Name)
	}

	// 5. Create session
	sessionService := session.InMemoryService()
	runnr, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          targetAgent,
		SessionService: sessionService,
		MemoryService:  memory.InMemoryService(), // Enables automatic memory persistence
	})
	if err != nil {
		log.Printf("Failed to create runner: %v", err)
		return nil, fmt.Errorf("failed to create runner: %w", err)
	}

	// 6. Init telemetry
	cfgTelemetry := loadConfig()
	telemetry, err := observability.NewTelemetry(ctx, cfgTelemetry)
	if err != nil {
		log.Printf("Failed to initialize telemetry: %v", err)
	}

	// tạo context mới để shutdown thay vì shutdown trực tiếp để tránh cancel context chính của app
	defer func() {
		shudownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(shudownCtx); err != nil {
			log.Printf("Error shutting down telemetry: %v", err)
		}
	}()
	log.Println("OpenTelemetry initialized successfully")

	return &AgentServer{Runner: runnr, SessionService: sessionService, Config: appCfg, Telemetry: telemetry}, nil
}

func (s *AgentServer) HandlerChat(c *gin.Context) {
	var r ChatRequest
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request": err.Error()})
		return
	}

	if r.SessionID == "" {
		sessionID := uuid.NewString()
		_, err := s.SessionService.Create(c.Request.Context(), &session.CreateRequest{
			UserID:    userID,
			SessionID: sessionID,
			AppName:   appName,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Failed to create session ": err.Error()})
			log.Printf("Failed to create session: %v", err)
			return
		}
		r.SessionID = sessionID
	}

	userMsg := genai.NewContentFromText(r.Message, genai.RoleUser)

	ctxOtel, span := s.Telemetry.Tracer.Start(c.Request.Context(), "agent.request")
	defer span.End()

	finalResponse := ""
	for event, err := range s.Runner.Run(ctxOtel, userID, r.SessionID, userMsg, agent.RunConfig{}) {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Failed agent response": err.Error()})
			log.Printf("Run error: %v", err)
		}
		if event.Content != nil && len(event.Content.Parts) > 0 {
			finalResponse += event.Content.Parts[0].Text
		}
	}

	c.JSON(http.StatusOK, ChatResponse{SessionID: r.SessionID, Message: finalResponse})
}
