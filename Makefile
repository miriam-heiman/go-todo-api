.PHONY: help build-lambda deploy-lambda test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-lambda: ## Build Lambda function for deployment
	@echo "Building Lambda function for ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap cmd/lambda/main.go
	@echo "✅ Lambda binary built: bootstrap"
	@ls -lh bootstrap

build-lambda-amd64: ## Build Lambda function for AMD64 (Intel)
	@echo "Building Lambda function for AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap cmd/lambda/main.go
	@echo "✅ Lambda binary built: bootstrap"
	@ls -lh bootstrap

deploy-lambda: build-lambda ## Deploy to AWS Lambda
	@echo "Deploying to AWS Lambda..."
	serverless deploy --verbose

deploy-lambda-dev: build-lambda ## Deploy to dev environment
	@echo "Deploying to dev environment..."
	serverless deploy --stage dev

deploy-lambda-prod: build-lambda ## Deploy to production
	@echo "Deploying to production..."
	serverless deploy --stage prod

remove-lambda: ## Remove Lambda deployment
	@echo "Removing Lambda deployment..."
	serverless remove

test: ## Run tests
	@echo "Running tests..."
	go test ./... -v -cover

test-lambda-local: build-lambda ## Test Lambda locally
	@echo "Testing Lambda locally..."
	serverless offline start

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f bootstrap
	rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

# Local development
run-local: ## Run API locally (HTTP mode)
	@echo "Starting API in HTTP mode..."
	go run cmd/api/main.go

# Docker commands
docker-build: ## Build Docker image
	docker-compose build api

docker-up: ## Start Docker stack
	docker-compose up -d

docker-down: ## Stop Docker stack
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f api
