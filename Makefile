BINARY_NAME := gh-purview
WAYBAR_BINARY := gh-purview-waybar
BUILD_DIR := bin
MAIN_PATH := ./cmd/app
WAYBAR_PATH := ./cmd/waybar

GO := go
GOFLAGS := -trimpath
LDFLAGS := -ldflags "-s -w"

.PHONY: all build build-waybar clean fmt vet test install install-waybar uninstall help

all: build

build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-waybar:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(WAYBAR_BINARY) $(WAYBAR_PATH)

install-waybar:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(shell go env GOPATH)/bin/$(WAYBAR_BINARY) $(WAYBAR_PATH)
	@echo "Installed $(WAYBAR_BINARY) to $(shell go env GOPATH)/bin"

clean:
	rm -rf $(BUILD_DIR)

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test -v -race ./...

install: build
	@mkdir -p $(HOME)/.local/share/gh/extensions/$(BINARY_NAME)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/share/gh/extensions/$(BINARY_NAME)/$(BINARY_NAME)
	@echo "Installed gh extension: $(BINARY_NAME)"

uninstall:
	rm -rf $(HOME)/.local/share/gh/extensions/$(BINARY_NAME)
	@echo "Uninstalled gh extension: $(BINARY_NAME)"

deps:
	$(GO) mod download
	$(GO) mod tidy

help:
	@echo "Available targets:"
	@echo "  build           - Build TUI for current platform"
	@echo "  build-waybar    - Build Waybar module for current platform"
	@echo "  clean           - Remove build artifacts"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  test            - Run tests"
	@echo "  deps            - Download and tidy dependencies"
	@echo "  install         - Install TUI as gh extension locally"
	@echo "  install-waybar  - Install Waybar module to GOBIN"
	@echo "  uninstall       - Remove local gh extension"
