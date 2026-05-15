package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LLMModel interface {
	Chat(ctx context.Context, prompt string) (string, error)
	StreamChat(ctx context.Context, prompt string, callback func(string)) (string, error)
}

type OpenAICompatibleLLM struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

func NewLLM(cfg OpenAICompatibleConfig) (LLMModel, error) {
	return &OpenAICompatibleLLM{
		BaseURL:    strings.TrimSuffix(cfg.BaseURL, "/"),
		APIKey:     cfg.APIKey,
		Model:      cfg.Model,
		HTTPClient: http.DefaultClient,
	}, nil
}

func (l *OpenAICompatibleLLM) Chat(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]any{
		"model": l.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", l.BaseURL+"/chat/completions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if l.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+l.APIKey)
	}

	resp, err := l.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llm api error: %s", string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from llm")
	}

	return result.Choices[0].Message.Content, nil
}

func (l *OpenAICompatibleLLM) StreamChat(ctx context.Context, prompt string, callback func(string)) (string, error) {
	reqBody := map[string]any{
		"model": l.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0,
		"stream":      true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", l.BaseURL+"/chat/completions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if l.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+l.APIKey)
	}

	resp, err := l.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llm api error: %s", string(body))
	}

	fullResp := ""
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fullResp, err
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			fullResp += content
			if callback != nil {
				callback(content)
			}
		}
	}

	return fullResp, nil
}
