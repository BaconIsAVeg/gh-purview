# AGENTS.md

Guidelines for AI coding agents working in this repository.

## Project Overview

`ghr` is a terminal UI (TUI) for browsing GitHub pull requests, built with the Charmbracelet ecosystem (Bubble Tea, Lipgloss, Bubbles).

## Build Gotcha

**IMPORTANT**: Always use `make build` or specify the output path with `-o bin/ghr`. Never run `go build ./cmd/app` without an output path, as it creates an executable in the current directory.

## Project Structure

```
ghr/
├── cmd/app/main.go          # Entry point
├── internal/
│   ├── types/               # Domain types (PR, PRStatus)
│   └── ui/
│       ├── model/           # Main Bubble Tea model (layout orchestration)
│       ├── header/          # Top bar component (filter, PR count)
│       ├── prlist/          # PR list component (two-line items)
│       ├── preview/         # Diff preview panel
│       ├── statusbar/       # Keybind hints bar
│       └── styles/          # Lipgloss color palette (centralized)
```

## Architecture Patterns

### Styling

- **All styles are centralized** in the `Palette` struct (`internal/ui/styles/styles.go`)
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
