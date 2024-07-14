# Makefile for Go projects

# Use PHONY to declare targets that don't represent files
.PHONY: build install tidy update test coverage clean

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=myapp
COVERAGE_FILE=coverage.out

.DEFAULT_GOAL := help

build:
	@echo "Building..."
	@$(GOBUILD) -o $(BINARY_NAME) -v

install:
	@echo "Installing..."
	@$(GOCMD) install

tidy:
	@echo "Tidying modules..."
	@$(GOCMD) mod tidy

update:
	@echo "Updating dependencies..."
	@$(GOGET) -u ./...
	@$(GOCMD) mod tidy

test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -covermode=atomic -coverprofile=$(COVERAGE_FILE) ./...
	@$(GOCMD) tool cover -html=$(COVERAGE_FILE)

clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE)

help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  install   - Install the application"
	@echo "  tidy      - Tidy Go modules"
	@echo "  update    - Update dependencies"
	@echo "  test      - Run tests"
	@echo "  coverage  - Run tests with coverage"
	@echo "  clean     - Clean build artifacts"