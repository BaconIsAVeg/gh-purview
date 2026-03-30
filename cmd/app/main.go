package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anomaly/ghr/internal/debug"
	"github.com/anomaly/ghr/internal/github"
	"github.com/anomaly/ghr/internal/ui/model"
	tea "github.com/charmbracelet/bubbletea"
)

var pageSize int

func init() {
	flag.IntVar(&pageSize, "pageSize", 25, "Number of PRs to fetch per request")
}

func main() {
	flag.Parse()

	if err := debug.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: debug init failed: %v\n", err)
	}
	defer debug.Close()

	ghClient, err := github.NewClient(pageSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing GitHub client: %v\n", err)
		os.Exit(1)
	}
	defer ghClient.Close()

	p := tea.NewProgram(model.New(ghClient), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
