ifneq (,$(wildcard .env))
  include .env
  export
endif

CMD_DIR := cmd/app


clean: ## Remove build artifacts and cache
	@echo "🧹 Cleaning up..."
	@rm -rf bin/ *.out *.exe *.test
	go clean


# Run the app
run: ## Run the app
	@echo "🚀 Running app:"
	go run $(CMD_DIR)/main.go


# --- Database Migrations ---
migrate-up: ## Run all migrations
	@echo "Running migrations..."
	go run cmd/migrate/main.go up

migrate-down: ## Rollback all migrations
	@echo "Rolling back migrations..."
	go run cmd/migrate/main.go down

migrate-version: ## Show current migration version
	go run cmd/migrate/main.go version

migrate-force: ## Force migration version (Usage: make migrate-force VERSION=1)
	go run cmd/migrate/main.go force $(VERSION)

migrate-steps: ## Run n migration steps (Usage: make migrate-steps STEPS=1)
	go run cmd/migrate/main.go steps $(STEPS)


# --- Tidy go.mod ---
tidy: ## Tidy go.mod and go.sum
	@echo "🧹 Tidying go.mod and go.sum..."
	go mod tidy


# --- Documentation ---
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort


test: ## Run all tests
	go test ./... 


test-force: ## Run tests without caching
	go test -count=1 ./... 


test-race: ## Run tests with race condition detection
	go test -race ./...


test-ci: ## Run tests with both race detection and coverage (used in CI)
	go test -race -coverprofile=coverage.out ./... 
	go tool cover -func=coverage.out


test-function: ## Usage: make test TEST=TestGetAllJournals
	go test -v -run ^$(TEST)$$


# Run all tests in the project,
test-log: ## Run all tests in the project, including showing logs
	go test -v ./... 


.PHONY: test test-force test-function run tidy help clean test-log migrate-up migrate-down migrate-version migrate-force migrate-steps