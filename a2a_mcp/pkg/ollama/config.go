package ollama

type Config struct {
	BaseURL    string // Mặc định là http://localhost:11434
	LLMModel   string // qwen2.5:0.5b
	EmbedModel string // qwen3-embedding:0.6b
}

func DefaultConfig() Config {
	return Config{
		BaseURL:    "http://localhost:11434",
		LLMModel:   "qwen2.5:0.5b",
		EmbedModel: "qwen3-embedding:0.6b",
	}
}
