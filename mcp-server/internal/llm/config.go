package llm

import "mcp-server/internal/utils"

func NewOpenAIEmbeddingConfig() OpenAICompatibleConfig {
	return OpenAICompatibleConfig{
		BaseURL:   utils.GetEnvString("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		Model:     utils.GetEnvString("OPENAI_EMBEDDING_MODEL", "text-embedding-3-large"),
		Dimension: utils.GetEnvInt("OPENAI_EMBEDDING_DIM", 1024),
		APIKey:    utils.GetEnvString("OPENAI_API_KEY", ""),
	}
}

func NewOpenAILLMConfig() OpenAICompatibleConfig {
	return OpenAICompatibleConfig{
		BaseURL: utils.GetEnvString("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		Model:   utils.GetEnvString("OPENAI_LLM_MODEL", "gpt-4o-mini"),
		APIKey:  utils.GetEnvString("OPENAI_API_KEY", ""),
	}
}

func NewOllamaEmbeddingConfig() OpenAICompatibleConfig {
	return OpenAICompatibleConfig{
		BaseURL:   utils.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
		Model:     utils.GetEnvString("OLLAMA_EMBEDDING_MODEL", "qwen3-embedding:0.6b"),
		Dimension: utils.GetEnvInt("OLLAMA_EMBEDDING_DIM", 1024),
	}
}

func NewOllamaLLMConfig() OpenAICompatibleConfig {
	return OpenAICompatibleConfig{
		BaseURL: utils.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
		Model:   utils.GetEnvString("OLLAM_LLM_MODEL", "qwen3.5:0.8b"),
	}
}
