# AGENTS.md

Guidelines for AI coding agents working in this repository.

## Project Overview

`gh-purview` is a terminal UI (TUI) for browsing GitHub pull requests, built with the Charmbracelet ecosystem (Bubble Tea, Lipgloss, Bubbles).

## Build Commands

This project uses [just](https://github.com/casey/just) as its task runner. Common recipes:

```bash
just build              # Build TUI binary to bin/gh-purview
just build-waybar       # Build Waybar module to bin/gh-purview-waybar
just test               # Run all tests with race detection
just fmt                # Format all Go code
just vet                # Run go vet
just deps               # Download and tidy dependencies
just install            # Install as gh extension locally
just install-waybar     # Install Waybar module to GOPATH/bin
just clean              # Remove build artifacts
just uninstall          # Remove local gh extension
```

## Build Gotcha

**IMPORTANT**: Always use `just build` or specify the output path with `-o bin/gh-purview`. Never run `go build ./cmd/app` without an output path, as it creates an executable in the current directory.

## Project Structure

```
gh-purview/
├── cmd/
│   ├── app/main.go           # TUI entry point
│   └── waybar/main.go        # Waybar module entry point
├── internal/
│   ├── types/                # Domain types (PR, PRStatus)
│   ├── github/               # GitHub API client (REST + GraphQL)
│   └── ui/
│       ├── model/            # Main Bubble Tea model (layout orchestration)
│       ├── appstyles/        # Lipgloss color palette (centralized)
│       ├── prlist/           # PR list component (two-line items)
│       ├── preview/          # Diff preview panel
│       └── statusbar/        # Keybind hints bar
```

## Architecture Patterns

### Styling

- **All styles are centralized** in the `Palette` struct (`internal/ui/appstyles/appstyles.go`)
- Pass palette pointer to child components via constructor
- Use 256-color codes for terminal compatibility
- When rendering styled text with backgrounds, ensure spaces are included in the rendered string to avoid transparent gaps

### Layout

- Use `lipgloss.JoinVertical` and `lipgloss.JoinHorizontal` for compositing
- Account for borders when calculating widths (subtract 2 for left/right borders)
- Call `updateLayout()` after window resize or panel toggle
- Use `EnsureCursorVisible()` after height changes to keep selection in view

### Bubble Tea Conventions

- Use import alias: `tea "github.com/charmbracelet/bubbletea"`
- Message types use `XxxMsg` suffix (e.g., `prsLoadedMsg`)
