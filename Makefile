BINARY_NAME := gh-purview
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := bin
MAIN_PATH := ./cmd/app

GO := go
GOFLAGS := -trimpath

LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: all build build-linux build-windows clean fmt vet test help

all: build

build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-linux: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

build-windows: clean
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)

build-all: build-linux build-windows

clean:
	rm -rf $(BUILD_DIR)

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test -v -race ./...

deps:
	$(GO) mod download
	$(GO) mod tidy

help:
	@echo "Available targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-linux   - Build static Linux binaries (amd64, arm64)"
	@echo "  build-windows - Build Windows executables (amd64, arm64)"
	@echo "  build-all     - Build all platform binaries"
	@echo "  clean         - Remove build artifacts"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  test          - Run tests"
	@echo "  deps          - Download and tidy dependencies"
