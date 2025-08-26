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

.PHONY: docker-start-kafka
docker-start-kafka: # Runs docker image
	@echo "$(BLUE) Starting docker container $(BLOCK_SCANNER_APP_NAME)...$(COLOR_END)"
	@docker-compose up -d kafka

.PHONY: run-app
run-app: docker-start-kafka # Runs application and wait for kafka
	@echo "$(BLUE) Waiting for Kafka to be ready...$(COLOR_END)"
	# Wait until Kafka port 9092 is open
	@sleep 5 #
	@while ! docker exec kafka sh -c "kafka-broker-api-versions --bootstrap-server kafka:9092" >/dev/null 2>&1; do \
		echo "Waiting for Kafka to be ready..."; \
		sleep 2; \
	done
	@echo "$(BLUE) Kafka is ready, starting application...$(COLOR_END)"
	@./bin/${BLOCK_SCANNER_APP_NAME}

.PHONY: docker-build
docker-build: # Builds docker image
	@echo "$(BLUE) Building docker image...$(COLOR_END)"
	@docker build --no-cache -t $(BLOCK_SCANNER_APP_NAME) .

.PHONY: docker-start
docker-start: # Runs docker image
	@echo "$(BLUE) Starting docker container $(BLOCK_SCANNER_APP_NAME)...$(COLOR_END)"
	@docker-compose up -d

.PHONY: docker-stop
docker-stop: # Runs docker image
	@echo "$(BLUE) Stopping and removing docker container $(BLOCK_SCANNER_APP_NAME)...$(COLOR_END)"
	@docker-compose down -v

lint: # Runs golangci-lint on the repo
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --timeout 5m

format: # Runs gofmt on the repo
	gofmt -s -w .

test: # Runs all the unit tests
	@echo "Test packages"
	go test -race -shuffle=on -coverprofile=coverage.out -cover $(PKGS)

.PHONY: help
help: # Show help for each of the Makefile recipes
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "$(GREEN)$$(echo $$l | cut -f 1 -d':')$(COLOR_END):$$(echo $$l | cut -f 2- -d'#')\n"; done
