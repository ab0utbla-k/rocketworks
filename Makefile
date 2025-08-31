BREW_PREFIX := $(shell brew --prefix)
BIN ?= $(BREW_PREFIX)/bin

## Tool Binaries
GOLANGCI_LINT ?= $(BIN)/golangci-lint

# ====================================================================================
# Colors

BLUE         := $(shell printf "\033[34m")
YELLOW       := $(shell printf "\033[33m")
RED          := $(shell printf "\033[31m")
GREEN        := $(shell printf "\033[32m")
CNone        := $(shell printf "\033[0m")

# ====================================================================================
# Logger

TIME_LONG	= `date +%Y-%m-%d' '%H:%M:%S`
TIME_SHORT	= `date +%H:%M:%S`
TIME		= $(TIME_SHORT)

INFO	= echo ${TIME} ${BLUE}[ .. ]${CNone}
WARN	= echo ${TIME} ${YELLOW}[WARN]${CNone}
ERR		= echo ${TIME} ${RED}[FAIL]${CNone}
OK		= echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL	= (echo ${TIME} ${RED}[FAIL]${CNone} && false)

# ====================================================================================
# OpenAPI Code Generation

.PHONY: generate
generate: ## Generate server code from OpenAPI spec
	@$(INFO) Bundling OpenAPI spec
	@redocly bundle order.openapi.yaml -o order.bundled.yaml
	@$(INFO) Generating Go server code
	@openapi-generator generate -i order.bundled.yaml -g go-server -o ./generated/server
	@$(OK) Generated server code from OpenAPI spec

.PHONY: clean-generated
clean-generated: ## Clean generated files
	@$(INFO) Cleaning generated files
	@rm -rf ./generated/
	@rm -f *.bundled.yaml *.bundled.yml
	@$(OK) Cleaned generated files

# ====================================================================================
# Golang
.PHONY: lint
lint: ## Run golangci-lint
	@if ! $(GOLANGCI_LINT) run; then \
		echo -e "\033[0;33mgolangci-lint failed: some checks can be fixed with \`\033[0;32mmake fmt\033[0m\033[0;33m\`\033[0m"; \
		exit 1; \
	fi
	@$(OK) Finished linting

.PHONY: fmt
fmt: ## Ensure consistent code style
	@go mod tidy
	@go fmt ./...
	@$(GOLANGCI_LINT) run --fix
	@$(OK) Ensured consistent code style

.PHONY: build
build: generate ## Build the order service binary
	@$(INFO) Building order service
	@go build -o bin/order ./cmd/order
	@$(OK) Built order service binary

.PHONY: test
test: ## Run tests
	@$(INFO) Running tests
	@go test ./...
	@$(OK) Tests passed

.PHONY: run
run: generate ## Run the order service
	@$(INFO) Starting order service
	@go run ./cmd/order
