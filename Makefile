BINARY_NAME := gh-purview
BUILD_DIR := bin
MAIN_PATH := ./cmd/app

GO := go
GOFLAGS := -trimpath
LDFLAGS := -ldflags "-s -w"

.PHONY: all build clean fmt vet test install uninstall help

all: build

build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

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
	@echo "  build     - Build for current platform"
	@echo "  clean     - Remove build artifacts"
	@echo "  fmt       - Format code"
	@echo "  vet       - Run go vet"
	@echo "  test      - Run tests"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  install   - Install as gh extension locally"
	@echo "  uninstall - Remove local gh extension"
