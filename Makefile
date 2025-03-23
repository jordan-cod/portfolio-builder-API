# Variáveis
APP_NAME=portfolio-backend
BIN_DIR=./bin
BIN_FILE=$(BIN_DIR)/$(APP_NAME)

# Comandos
.PHONY: build run dev clean fmt lint test

build:
	go build -o $(BIN_FILE) ./cmd

run: build
	$(BIN_FILE)

dev:
	air

clean:
	rm -rf $(BIN_DIR)

fmt:
	go fmt ./...
	goimports -w .

lint:
	golangci-lint run

test:
	go test ./...

