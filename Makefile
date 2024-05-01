.DEFAULT_GOAL := default

build:
	@echo "Building..."
	@go build

install:
	@go install

tidy:
	@echo "Making mod tidy"
	@go mod tidy

update:
	@echo "Updating..."
	@go get -u ./...
	@go mod tidy

coverage:
	@echo "Testing..."
	@go test -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

default: tidy build
