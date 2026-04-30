package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Instruction string   `yaml:"instruction"`
	Tools       []string `yaml:"tools"`
}

type AppConfig struct {
	Agents    map[string]AgentConfig `yaml:"agents"`
	McpServer string                 `yaml:"mcp_server"`
	Models    map[string]string      `yaml:"models"`
}

func LoadAgentConfig(filePath string, agentName string) (*AgentConfig, error) {
	// 1. Read File yaml
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 2. Parse yaml
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	agent, ok := config.Agents[agentName]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", agentName)
	}
	return &agent, nil
}

func LoadAppConfig(filePath string) (*AppConfig, error) {
	// 1. Read File yaml
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 2. Parse yaml
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
