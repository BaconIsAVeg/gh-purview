package github

import (
	"context"
	"fmt"

	"github.com/anomaly/ghr/internal/types"
	"github.com/google/go-github/v82/github"
)

const maxDiffSize = 100 * 1024 // 100KB limit for diffs

type DiffResult struct {
	Content      string
	Truncated    bool
	Additions    int
	Deletions    int
	ChangedFiles int
}

func (c *Client) FetchPRDiff(ctx context.Context, pr *types.PR) (*DiffResult, error) {
	if pr == nil {
		return nil, fmt.Errorf("PR is nil")
	}

	prDetail, _, err := c.client.PullRequests.Get(ctx, pr.Org, pr.Repo, pr.Number)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR details: %w", err)
	}

	result := &DiffResult{
		Additions:    prDetail.GetAdditions(),
		Deletions:    prDetail.GetDeletions(),
		ChangedFiles: prDetail.GetChangedFiles(),
	}

	if result.ChangedFiles > 100 {
		files, _, err := c.client.PullRequests.ListFiles(ctx, pr.Org, pr.Repo, pr.Number, &github.ListOptions{PerPage: 100})
		if err != nil {
			return nil, fmt.Errorf("failed to list PR files: %w", err)
		}
		result.Content = formatFileList(files, result.Additions, result.Deletions, result.ChangedFiles)
		result.Truncated = true
		return result, nil
	}

	diff, _, err := c.client.PullRequests.GetRaw(ctx, pr.Org, pr.Repo, pr.Number, github.RawOptions{Type: github.Diff})
	if err != nil {
		return nil, fmt.Errorf("failed to get PR diff: %w", err)
	}

	if len(diff) > maxDiffSize {
		result.Content = diff[:maxDiffSize] + "\n\n... (diff truncated - too large)"
		result.Truncated = true
		return result, nil
	}

	result.Content = diff
	return result, nil
}

func formatFileList(files []*github.CommitFile, additions, deletions, changedFiles int) string {
	var b string
	b = fmt.Sprintf("Large PR: %d files changed, only showing file list.\n\n", changedFiles)
	b += "Files changed:\n"
	for i, f := range files {
		if i >= 50 {
			b += fmt.Sprintf("  ... and %d more files\n", changedFiles-50)
			break
		}
		status := f.GetStatus()
		filename := f.GetFilename()
		adds := f.GetAdditions()
		dels := f.GetDeletions()
		b += fmt.Sprintf("  %s %s (+%d/-%d)\n", statusIcon(status), filename, adds, dels)
	}
	return b
}

func statusIcon(status string) string {
	switch status {
	case "added":
		return "A"
	case "removed":
		return "D"
	case "modified":
		return "M"
	case "renamed":
		return "R"
	default:
		return "?"
	}
}
