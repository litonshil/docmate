PROJECT_NAME := docmate
PKG_LIST := $(shell go list ${PROJECT_NAME}/... | grep -v /vendor/)

.PHONY: all dep build clean test lint help development

all: build ## Build the project

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

########################
### DEVELOP and TEST ###
########################

development: ## Set up development environment
	# Booting up dependency containers
	@docker-compose up -d consul db redis

	# Wait for consul container to be ready
	@while ! curl --request GET -sL --url 'http://localhost:8500/' > /dev/null 2>&1; do printf .; sleep 1; done

	# Setting KV, dependency of app
	@curl --request PUT --data-binary @config.local.json http://localhost:8500/v1/kv/${PROJECT_NAME} || { echo "Failed to set KV"; exit 1; }

	# Building docmate
	@docker-compose up --build ${PROJECT_NAME}

lint: ## Run golangci-lint (v2)
	@which golangci-lint > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.5.0
	@$(shell go env GOPATH)/bin/golangci-lint run ./...

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)
	@docker-compose down

build: ## Build the project
	# Build commands here if needed
	@echo "Build completed"
