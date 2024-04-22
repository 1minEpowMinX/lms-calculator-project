run-orchestrator:
	go run cmd/orchestrator/main.go

run-agent:
	go run cmd/agent/main.go

run-project: run-orchestrator run-agent