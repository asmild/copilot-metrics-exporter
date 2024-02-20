BIN := copilot-metrics-exporter
SRC := ./cmd/copilot-exporter

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
