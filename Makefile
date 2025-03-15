# Define the binary name
BINARY_NAME=./out/aslan-words
BINARY_DEMO_NAME=./demo/aslan-words

# Define the main package
MAIN_PACKAGE=./cmd/generate

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

# Build binary to generate the demo.gif inside the vhs docker container
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

save-demo-gif: build-demo
	docker run --rm -v ${PWD}/demo:/vhs ghcr.io/charmbracelet/vhs demo.tape

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
