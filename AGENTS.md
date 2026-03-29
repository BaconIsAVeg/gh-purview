# AGENTS.md

Guidelines for AI coding agents working in this repository.

## Project Overview

`ghr` is a terminal UI (TUI) for browsing GitHub pull requests, built with the Charmbracelet ecosystem (Bubble Tea, Lipgloss, Bubbles).

## Build Commands

**IMPORTANT**: Always use `make build` or specify the output path with `-o bin/ghr`. Never run `go build ./cmd/app` without an output path, as it creates an executable in the current directory.

```bash
# Build for current platform (recommended)
make build
# Or with explicit output path:
go build -trimpath -ldflags "-s -w" -o bin/ghr ./cmd/app

# Build static Linux binaries
make build-linux

# Build Windows executables
make build-windows

# Build all platforms
make build-all

# Clean build artifacts
make clean
```

## Test Commands

```bash
# Run all tests
make test
# Or: go test -v -race ./...

# Run a single test
go test -v -race -run TestFunctionName ./path/to/package

# Run tests in a specific file
go test -v -race -run TestFunctionName ./path/to/package/file_test.go

# Run tests matching a pattern
go test -v -race -run "TestPattern.*" ./...
```

## Lint/Format Commands

```bash
# Format code
make fmt
# Or: go fmt ./...

# Run go vet
make vet
# Or: go vet ./...

# Tidy dependencies
make deps
# Or: go mod tidy
```

## Project Structure

```
ghr/
├── cmd/app/main.go          # Entry point
├── internal/
│   ├── types/               # Domain types (PR, PRStatus)
│   └── ui/
│       ├── model/           # Main Bubble Tea model (layout orchestration)
│       ├── header/          # Top bar component (filter, PR count)
│       ├── prlist/           # PR list component (two-line items)
│       ├── preview/          # Diff preview panel
│       ├── statusbar/        # Keybind hints bar
│       └── styles/           # Lipgloss color palette (centralized)
├── go.mod
├── Makefile
└── .gitignore
```

## Code Style Guidelines

### Imports

- Group imports: standard library, external packages, internal packages
- Use import aliases for `bubbletea`: `tea "github.com/charmbracelet/bubbletea"`
- Example:

  ```go
  import (
      "fmt"
      "strings"

      "github.com/anomaly/ghr/internal/types"
      "github.com/anomaly/ghr/internal/ui/styles"
      tea "github.com/charmbracelet/bubbletea"
      "github.com/charmbracelet/lipgloss"
  )
  ```

### Formatting

- Use `gofmt` (or `go fmt`) for formatting
- No trailing whitespace
- Tabs for indentation (Go standard)

### Types

- Define custom types for domain concepts (e.g., `PRStatus`, `PR`)
- Use string-backed enums for status types:

  ```go
  type PRStatus string

  const (
      StatusOpen   PRStatus = "open"
      StatusClosed PRStatus = "closed"
      StatusMerged PRStatus = "merged"
  )
  ```

### Naming Conventions

- **Packages**: lowercase, single word when possible (`types`, `prlist`, `preview`)
- **Types**: PascalCase (`Model`, `Palette`, `KeyBinding`)
- **Functions/Methods**: PascalCase for exported, camelCase for private
- **Message types**: `XxxMsg` suffix (e.g., `prsLoadedMsg`)

### Bubble Tea Patterns

- Models implement `Init()`, `Update()`, and `View()` methods
- Use custom message types for internal communication:
  ```go
  type prsLoadedMsg struct {
      prs []types.PR
  }
  ```
- Return `tea.Cmd` functions for async operations
- Use `tea.Batch()` to combine multiple commands
- Handle `tea.WindowSizeMsg` for responsive layouts

### Styling with Lipgloss

- **Centralize all styles** in `Palette` struct (`internal/ui/styles/styles.go`)
- Pass palette pointer to child components via constructor
- Define reusable color variables at function level:
  ```go
  primaryBg := lipgloss.Color("97")
  secondaryBg := lipgloss.Color("234")
  ```
- Use 256-color codes for terminal compatibility
- When rendering styled text with backgrounds, ensure spaces are included in the rendered string to avoid transparent gaps

### Layout Patterns

- Use `lipgloss.JoinVertical` and `lipgloss.JoinHorizontal` for compositing
- Account for borders when calculating widths (subtract 2 for left/right borders)
- Call `updateLayout()` after window resize or panel toggle
- Use `EnsureCursorVisible()` after height changes to keep selection in view

### Error Handling

- Handle errors explicitly; do not ignore them
- Use `fmt.Printf` for simple error output in TUI context
- Exit with `os.Exit(1)` for fatal errors in `main()`

## Key Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework (Elm Architecture)
- `github.com/charmbracelet/lipgloss` - Styling/layout engine
- `github.com/charmbracelet/bubbles` - Pre-built UI components (viewport)

## Key Bindings

- `j`/`k`: Navigate list
- `p`: Toggle preview panel
- `r`: Refresh
- `q` or `Ctrl+C`: Quit
