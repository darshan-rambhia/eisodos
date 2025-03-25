#! make
.DEFAULT_GOAL := help

# Project variables
BINARY_NAME := eisodos
BUILD_DIR := .out
BINARY_PATH := $(BUILD_DIR)/bin/$(BINARY_NAME)
COVERAGE_DIR := $(BUILD_DIR)/coverage
TEST_DIR := $(BUILD_DIR)/test
APP_BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S)
APP_GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "dev")
SERVICE_PORT := 8080

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOCLEAN := $(GOCMD) clean
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -ldflags "-w -s -X main.BuildDate=$(APP_BUILD_DATE) -X main.GitCommit=$(APP_GIT_COMMIT)"

# Test flags
TESTFLAGS := -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out

# Setup environment variables
ifneq (,$(wildcard $(shell pwd)/.env))
	include .env
	export
endif

help: ## Show this list of commands
	@printf "This Makefile handles eisodos load balancer development and deployment with following commands:\n"
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

init: ## Initialize the project and download dependencies
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy

build: init ## Build the binary for local runs
	mkdir -p $(BUILD_DIR)/bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/eisodos

test: ## Run tests with gotestsum (includes JUnit and Testdox reports)
	mkdir -p $(TEST_DIR)
	$(GOCMD) run buildscripts/scripts/test.go

clean: ## Clean build artifacts and temporary files
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.html

run: build ## Run the application locally
	$(BINARY_PATH)

coverage: ## Generate HTML coverage report
	mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TESTFLAGS) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

lint: ## Run golangci-lint
	golangci-lint run

install-tools: ## Install development tools (golangci-lint, gotestsum)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install gotest.tools/gotestsum@latest

.PHONY: docker_test
.ONESHELL: docker_test
docker_test: ## Build and run docker test image locally
	@printf "\nBuilding docker image $(BINARY_NAME)-test\n"
	docker compose -f buildscripts/docker/docker-test.yaml build
	@printf "\nRunning docker image $(BINARY_NAME)-test\n"
	docker compose -f buildscripts/docker/docker-test.yaml up
	@printf "\nStopping docker image $(BINARY_NAME)-test\n"
	docker compose -f buildscripts/docker/docker-test.yaml down --volumes
	@printf "\nCleanup complete!\n"

.PHONY: docker_build
docker_build: ## Build docker image
	@printf "Building docker image $(BINARY_NAME)\n"
	docker compose -f buildscripts/docker/docker-run.yaml build
	@printf "\n$(BINARY_NAME) image is now available\n"

.PHONY: docker_run
.ONESHELL: docker_run
docker_run: ## Run docker container
	@printf "\nRunning docker container $(BINARY_NAME)\n"
	docker compose -f buildscripts/docker/docker-run.yaml up --detach
	@printf "\n$(BINARY_NAME) is up and running on port $(SERVICE_PORT)\n"

docker_cleanup: ## Cleanup docker resources
	docker compose -f buildscripts/docker/docker-run.yaml down --volumes

healthcheck: ## Run healthcheck on the running application
	@printf "Running healthcheck on the application\n"
	$(eval RESPONSE=$(shell curl --silent "http://localhost:$(SERVICE_PORT)/health"))
	if [ "$$RESPONSE" == "ok" ]; then\
		printf "Healthcheck succeeded\n";\
	else\
		printf "Healthcheck failed, response: $(RESPONSE)\n";\
	fi

docker_run_test: ## Build, run, healthcheck, and cleanup docker container
	make docker_build
	make docker_run
	@printf "\n"
	make healthcheck
	@printf "\n"
	make docker_cleanup

$(VERBOSE).SILENT: 