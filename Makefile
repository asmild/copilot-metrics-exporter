BIN := copilot-metrics-exporter
SRC := ./cmd/copilot-exporter
TAG := latest

.PHONY: run
run: build
	./$(BIN)

.PHONY: build
build:
	@echo "Building $(BIN)..."
	@go build -o $(BIN) $(SRC)/main.go

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BIN)

.PHONY: image
image:
	@echo "Building Docker image..."
	@docker build -t asmild/copilot-metrics-exporter:$(TAG) .