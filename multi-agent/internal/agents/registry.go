package agents

import (
	"fmt"
	"sync"

	"google.golang.org/adk/agent"
)

// Registry manages the lifecycle and discovery of agents.
type Registry struct {
	mu     sync.RWMutex
	agents map[string]agent.Agent
}

// NewRegistry creates a new agent registry.
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]agent.Agent),
	}
}

// Register adds an agent to the registry.
func (r *Registry) Register(name string, a agent.Agent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[name] = a
}

// Get retrieves an agent from the registry by name.
func (r *Registry) Get(name string) (agent.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.agents[name]
	if !ok {
		return nil, fmt.Errorf("agent %s not found in registry", name)
	}
	return a, nil
}

// ListNames returns a list of all registered agent names.
func (r *Registry) ListNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.agents))
	for name := range r.agents {
		names = append(names, name)
	}
	return names
}

// Get Agents
func (r *Registry) GetAgents() []agent.Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	agents := make([]agent.Agent, 0, len(r.agents))
	for name := range r.agents {
		agents = append(agents, r.agents[name])
	}
	return agents
}
