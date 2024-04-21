package main

import "github.com/1minepowminx/distributed_calculator/internal/grpc/agent"

// Agent - the calculating server, which parses expression ->
// calculates the answer -> returns result to Orchestrator

func main() {
	agent.RunAgentServer()
}
