package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BaconIsAVeg/gh-purview/internal/config"
	"github.com/BaconIsAVeg/gh-purview/internal/github"
	"github.com/BaconIsAVeg/gh-purview/internal/ui/model"
	"github.com/BaconIsAVeg/github-tuis/debug"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	pageSize   int
	filterName string
)

func init() {
	flag.IntVar(&pageSize, "pageSize", 25, "Number of PRs to fetch per request")
	flag.StringVar(&filterName, "filter", "", "Named filter from config to use at startup")
}

func main() {
	flag.Parse()

	if err := debug.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: debug init failed: %v\n", err)
	}
	defer debug.Close()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: config load failed: %v\n", err)
		cfg = &config.Config{}
	}
	last, err := config.LoadLast()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: resume state load failed: %v\n", err)
	}
	query, activeFilter, err := config.Resolve(cfg, last, filterName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ghClient, err := github.NewClient(pageSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing GitHub client: %v\n", err)
		os.Exit(1)
	}
	defer ghClient.Close()
	if query != "" {
		ghClient.SetQuery(query)
	}

	p := tea.NewProgram(model.New(ghClient, activeFilter), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
