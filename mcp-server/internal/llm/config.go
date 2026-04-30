package llm

type Provider string

const (
	ProviderOllama Provider = "ollama"
	ProviderOpenAI Provider = "openai"
	ProviderGroq   Provider = "groq"
)

type ProviderConfig struct {
	Provider Provider
	BaseURL  string
	APIKey   string
	Model    string
}

type Config struct {
	LLM   ProviderConfig
	Embed ProviderConfig
}

func DefaultConfig() Config {
	return Config{
		LLM: ProviderConfig{
			Provider: ProviderOllama,
			BaseURL:  "http://localhost:11434",
			Model:    "qwen2.5:0.5b",
		},
		Embed: ProviderConfig{
			Provider: ProviderOllama,
			BaseURL:  "http://localhost:11434",
			Model:    "qwen3-embedding:0.6b",
		},
	}
}


