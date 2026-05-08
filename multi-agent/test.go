package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"multi-agent/internal/agents"
	"multi-agent/internal/config"
	mymcp "multi-agent/internal/mcp"
	"multi-agent/internal/observability"
	"multi-agent/internal/provider/gemini"
	"multi-agent/internal/utils"

	"net/http"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool/toolconfirmation"
	"google.golang.org/genai"
)

// Agent (Start Span) -> RoundTrip (Inject) -> [Network] -> MCP Server (Extract) -> Tool Handler (Inject) -> [Network] -> Backend (Extract)
const (
	appName = "salon_chain"
	userID  = "demo_user"
)

type Config struct {
	Port          string
	Telemetry     observability.Config
	EnableTracing bool
	EnableMetrics bool
}

func loadConfig() Config {
	return Config{
		Telemetry: observability.Config{
			ServiceName:    utils.GetEnvString("OTEL_SERVICE", "agent-server"),
			ServiceVersion: utils.GetEnvString("OTEL_VERSION", "1.0.0"),
			Environment:    utils.GetEnvString("ENVIRONMENT", "development"),
			OTLPEndpoint:   utils.GetEnvString("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
			SamplingRate:   utils.GetEnvFloat("OTEL_TRACES_SAMPLER_ARG", 1.0),
			EnableTracing:  utils.GetEnvBool("OTEL_ENABLE_TRACING", true),
			EnableMetrics:  utils.GetEnvBool("OTEL_ENABLE_METRICS", true),
		},
	}
}

func main() {
	ctx := context.Background()

	// 1. Cấu hình cứng để test
	mcpServer := "http://localhost:8081/mcp"
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	path := "./config.yaml"
	agentCfg, err := config.LoadAgentConfig(path, "qa_agent")
	if err != nil {
		fmt.Println("failed to init config agent")
	}

	// 2. Init LLM
	llm, err := gemini.NewGeminiLLM(ctx, "gemini-2.5-flash-lite")
	if err != nil {
		log.Fatalf("Failed to create LLM: %v", err)
	}

	// 3. Init MCP Toolset với custom RoundTripper để chèn Auth Header
	mcpToken := os.Getenv("AUTH_TOKEN")

	// Tạo một http.Client tùy chỉnh để chèn header Authorization
	customClient := &http.Client{
		Transport: &headerTransport{
			base:   http.DefaultTransport,
			header: "Authorization",
			value:  "Bearer " + mcpToken,
		},
	}

	transport := &mcp.SSEClientTransport{
		Endpoint:   mcpServer,
		HTTPClient: customClient,
	}

	mcpToolset, err := mymcp.NewMCPTool(transport, agentCfg.AllowedTools, agentCfg.ApprovedTools)
	if err != nil {
		log.Fatalf("Failed to create MCP tools: %v", err)
	} else {
		log.Printf("Init mcp tool set successfull")
	}

	// 4. Create Single Agent
	testAgent, err := agents.NewSubAgent(ctx, agentCfg, llm, mcpToolset)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 5. Init Runner
	sessionService := session.InMemoryService()
	runnr, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          testAgent,
		SessionService: sessionService,
		MemoryService:  memory.InMemoryService(),
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	// 6. Create Session
	sessionID := uuid.NewString()
	_, err = sessionService.Create(ctx, &session.CreateRequest{
		UserID:    userID,
		SessionID: sessionID,
		AppName:   appName,
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// 7. Run Test Chat
	userInput := "Bên bạn có những chi nhánh nào?"
	fmt.Printf("--- User: %s ---\n", userInput)

	userMsg := genai.NewContentFromText(userInput, genai.RoleUser)

	// 8. Chạy trong một Span để truyền context đi
	cfg := loadConfig()
	telemetry, err := observability.NewTelemetry(ctx, cfg.Telemetry)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}

	// tạo context mới để shutdown thay vì shutdown trực tiếp để tránh cancel context chính của app
	defer func() {
		shudownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(shudownCtx); err != nil {
			log.Fatalf("Error shutting down telemetry: %v", err)
		}
	}()
	log.Println("OpenTelemetry initialized successfully")

	ctxOtel, span := telemetry.Tracer.Start(ctx, "agent.request")
	defer span.End()

	// Turn 1: chạy bình thường, capture confirmation event
	var pendingConfirmations map[string]toolconfirmation.ToolConfirmation
	var confirmationCallID string

	for event, err := range runnr.Run(ctxOtel, userID, sessionID, userMsg, agent.RunConfig{}) {
		if err != nil {
			log.Fatalf("Run error: %v", err)
		}

		fmt.Println(event.Actions)

		// In text bình thường
		if event.Content != nil && len(event.Content.Parts) > 0 {
			if event.Content.Parts[0].Text != "" {
				fmt.Println(event.Content.Parts[0].Text)
			}
			for _, part := range event.Content.Parts {
				// Nếu thấy Agent gọi hàm adk_request_confirmation
				if part.FunctionCall != nil && part.FunctionCall.Name == "adk_request_confirmation" {
					confirmationCallID = part.FunctionCall.ID // LẤY CÁI ID NÀY!
					log.Printf("[DEBUG] Thấy hàm duyệt ảo! ID: %s", confirmationCallID)
				}
			}
		}
		// Capture confirmation đang chờ
		if event.Actions.RequestedToolConfirmations != nil {
			pendingConfirmations = event.Actions.RequestedToolConfirmations
		}
	}
	fmt.Println("Pending confirm:", pendingConfirmations)

	// Turn 2: nếu có confirmation đang chờ → gửi approve response
	if pendingConfirmations != nil {
		var parts []*genai.Part
		for callID, conf := range pendingConfirmations {
			log.Printf("[DEBUG] Approving: callID=%s, hint=%s", callID, conf.Hint)

			if confirmationCallID != "" {
				// Dùng confirmationCallID thay vì cái ID trong map RequestedToolConfirmations
				parts = append(parts, &genai.Part{
					FunctionResponse: &genai.FunctionResponse{
						ID:   confirmationCallID, // PHẢI DÙNG ID NÀY
						Name: "adk_request_confirmation",
						Response: map[string]any{
							"confirmed": true,
							"hint":      conf.Hint,
							"payload":   conf.Payload,
						},
					},
				})
			}
		}

		approvalMsg := &genai.Content{
			Role:  string(genai.RoleUser),
			Parts: parts,
		}

		for event, err := range runnr.Run(ctxOtel, userID, sessionID, approvalMsg, agent.RunConfig{}) {
			if err != nil {
				log.Fatalf("Resume error: %v", err)
			}

			fmt.Println(event.Content)
			fmt.Println(event.Actions)
			if event.Content != nil && len(event.Content.Parts) > 0 {
				fmt.Print(event.Content.Parts[0].Text)
			}
		}
	}
}

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

// C:\Users\duc.phamtrinh\.gemini\antigravity\brain\78ea7052-df2b-49f3-ae9e-af0026772ee1\distributed_tracing_flow.md.resolved
