package graph

import "github.com/victor-nazario/kube-agent/internal/agent"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	agent agent.Agent
}
