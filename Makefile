BIN := copilot-metrics-exporter
SRC := ./cmd/copilot-exporter
TAG := latest

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

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

.PHONY: push
push:
	@echo "Pushing Docker image..."
	@docker push asmild/copilot-metrics-exporter:$(TAG)