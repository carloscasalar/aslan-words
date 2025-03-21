# Define the binary name
BINARY_NAME=./out/generate-word
BINARY_DEMO_NAME=./demo/generate-word

# Define the main package
MAIN_PACKAGE=./cmd/generate-word

# Go Bin
GO_BIN=$(shell which go)

# Ensure the out directory exists
out:
	mkdir -p out

install-tools: ## Install the required tools for go tools to develop
	@go install tool

build: out ## Build the Go application
	@go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

build-demo: out ## Build binary to generate-word the demo.gif inside the vhs docker container
	env GOOS=linux go build -o $(BINARY_DEMO_NAME) $(MAIN_PACKAGE)

run: build ## Run the Go application, you can use -- to pass arguments to the binary
	$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

test: ## Run tests
	@go test -v ./...

lint: ## Run the linter
	@revive -config .revive.toml -formatter friendly ./...

check: lint test ## Run all checks: lint and test

save-demo-gif: build-demo ## Generate a demo gif
	docker run --rm -v ${PWD}/demo:/vhs ghcr.io/charmbracelet/vhs demo.tape

publish-demo-gif: build-demo ## Publish the demo gif
	docker run --rm -v ${PWD}/demo:/vhs ghcr.io/charmbracelet/vhs demo.tape --publish

clean: ## Clean up build artifacts
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_DEMO_NAME)

# Default target
all: build

# Hack to make run proxy the arguments to the binary
%:
	@true

.PHONY: out build run test lint check clean all

# Show help
help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
