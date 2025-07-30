.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.SILENT: build
build: ## Build Webhook
	GOOS="linux" GOARCH="amd64" go build -o webhook-example ./cmd/main.go

.PHONY: start
start: ## Start the webhook server locally with default settings
	@echo "Starting webhook server with default configuration..."
	@echo "Server will be available at http://localhost:8080"
	@echo "Health check: http://localhost:8080/healthz"
	@echo "Webhook endpoint: http://localhost:8080/api/hook"
	@echo ""
	SHARED_SECRET=development-secret \
	OUTPUT_TYPE=stdout \
	PORT=8080 \
	go run cmd/main.go