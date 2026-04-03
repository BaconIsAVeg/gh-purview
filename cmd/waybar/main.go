package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/BaconIsAVeg/gh-purview/internal/github"
)

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

func main() {
	client, err := github.NewClient(100)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	_, total, err := client.FetchPRs(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching PRs: %v\n", err)
		os.Exit(1)
	}

	count := total
	var class string
	if count > 0 {
		class = "has-prs"
	} else {
		class = "no-prs"
	}

	output := WaybarOutput{
		Text:    fmt.Sprintf("%d", count),
		Tooltip: client.Query(),
		Class:   class,
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonOutput))
}
