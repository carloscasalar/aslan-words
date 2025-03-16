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

# Install the required tools for go generators
install-tools:
	@go install tool

# Build the Go application
build: out
	@go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Build binary to generate-word the demo.gif inside the vhs docker container
build-demo: out
	env GOOS=linux go build -o $(BINARY_DEMO_NAME) $(MAIN_PACKAGE)

# Run the Go application
run: build
	$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

# Run tests
test:
	@go test -v ./...

# Run the linter
lint:
	@revive -config .revive.toml -formatter friendly ./...

# Run all checks
check: lint test

# Generate a demo gif
save-demo-gif: build-demo
	docker run --rm -v ${PWD}/demo:/vhs ghcr.io/charmbracelet/vhs demo.tape

# Publish the demo gif
publish-demo-gif: build-demo
	docker run --rm -v ${PWD}/demo:/vhs ghcr.io/charmbracelet/vhs demo.tape --publish

# Clean up build artifacts
clean:
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
