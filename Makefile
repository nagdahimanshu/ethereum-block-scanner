BLOCK_SCANNER_APP_NAME = block-scanner
PKGS=$(shell go list ./... | grep -v "/vendor/")
BLUE = \033[1;34m
GREEN = \033[1;32m
COLOR_END = \033[0;39m
GO_VERSION = $(shell grep "^go " go.mod | head -n1 | cut -d " " -f2)

install-deps: # Installs the dependencies
	@echo "Installing necessary go dependencies for the api"
	go mod tidy

build: # Builds the application and create a binary at ./bin/
	@echo "$(BLUE)Building $(BLOCK_SCANNER_APP_NAME) binary...$(COLOR_END)"
	go build -a -o ./bin/$(BLOCK_SCANNER_APP_NAME) ./cmd/${BLOCK_SCANNER_APP_NAME}/...
	@echo "$(GREEN)Binary successfully built$(COLOR_END)"

lint: # Runs golangci-lint on the repo
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --timeout 5m

format: # Runs gofmt on the repo
	gofmt -s -w .

.PHONY: help
help: # Show help for each of the Makefile recipes
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "$(GREEN)$$(echo $$l | cut -f 1 -d':')$(COLOR_END):$$(echo $$l | cut -f 2- -d'#')\n"; done
