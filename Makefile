.PHONY: build test lint docker-build tidy help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Services
API_GATEWAY_DIR=services/api-gateway

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build all services"
	@echo "  test          Run tests for all services"
	@echo "  lint          Run linter"
	@echo "  docker-build  Build Docker images for all services"
	@echo "  tidy          Run go mod tidy for all services"
	@echo "  help          Show this help message"

build: build-api-gateway

build-api-gateway:
	cd $(API_GATEWAY_DIR) && $(GOBUILD) -v ./...

test: test-root test-api-gateway

test-root:
	$(GOTEST) -v ./...

test-api-gateway:
	cd $(API_GATEWAY_DIR) && $(GOTEST) -v ./...

lint:
	golangci-lint run ./...
	cd $(API_GATEWAY_DIR) && golangci-lint run ./...

docker-build: docker-build-api-gateway

docker-build-api-gateway:
	docker build -t api-gateway ./services/api-gateway

tidy:
	$(GOMOD) tidy
	cd $(API_GATEWAY_DIR) && $(GOMOD) tidy
