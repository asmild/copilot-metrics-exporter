BIN := copilot-metrics-exporter
SRC := ./cmd/copilot-exporter
TAG := latest
BUILD_DIR := build
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./... -coverprofile=coverage.out

.PHONY: run
run: build
	./$(BIN)

.PHONY: build
build:
	@echo "Building $(BIN) for current platform..."
	@go build $(LDFLAGS) -o $(BIN) $(SRC)/main.go

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BIN)
	@rm -rf $(BUILD_DIR)

# Cross-platform builds
.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "Building $(BIN) for Linux amd64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-linux-amd64 $(SRC)/main.go

.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "Building $(BIN) for Linux arm64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-linux-arm64 $(SRC)/main.go

.PHONY: build-darwin-amd64
build-darwin-amd64:
	@echo "Building $(BIN) for macOS amd64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-darwin-amd64 $(SRC)/main.go

.PHONY: build-darwin-arm64
build-darwin-arm64:
	@echo "Building $(BIN) for macOS arm64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-darwin-arm64 $(SRC)/main.go

.PHONY: build-windows-amd64
build-windows-amd64:
	@echo "Building $(BIN) for Windows amd64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-windows-amd64.exe $(SRC)/main.go

.PHONY: build-all
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
	@echo "All binaries built successfully in $(BUILD_DIR)/"

.PHONY: list-builds
list-builds:
	@echo "Available binaries:"
	@ls -lh $(BUILD_DIR)/ 2>/dev/null || echo "No builds found. Run 'make build-all' first."

# Docker targets (kept for backward compatibility)
.PHONY: image
image:
	@echo "Building Docker image..."
	@docker build -t asmild/copilot-metrics-exporter:$(TAG) .

.PHONY: push
push:
	@echo "Pushing Docker image..."
	@docker push asmild/copilot-metrics-exporter:$(TAG)