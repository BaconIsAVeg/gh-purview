package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/anomaly/ghr/internal/debug"
	"github.com/anomaly/ghr/internal/types"
	"github.com/google/go-github/v82/github"
)

func debugPrint(format string, args ...interface{}) {
	debug.Print(format, args...)
}

const prSearchQuery = "org:%s is:pr is:open -is:draft sort:updated-desc"

func (c *Client) FetchPRs(ctx context.Context) ([]types.PR, int, error) {
	query := fmt.Sprintf(prSearchQuery, c.org)
	debugPrint("Query: %s", query)

	perPage := c.pageSize
	if perPage > 100 {
		perPage = 100
	}

	var allPRs []types.PR
	var totalCount int
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: perPage},
	}

	page := 1
	for {
		debugPrint("Fetching page %d...", page)
		result, resp, err := c.client.Search.Issues(ctx, query, opts)
		if err != nil {
			debugPrint("Search error: %v", err)
			return nil, 0, fmt.Errorf("search failed: %w", err)
		}

		totalCount = result.GetTotal()
		debugPrint("Response: TotalCount=%d, ItemsInPage=%d, NextPage=%d, Rate.Remaining=%d",
			totalCount, len(result.Issues), resp.NextPage, resp.Rate.Remaining)
		debugPrint("IncompleteResults=%v", result.GetIncompleteResults())

		if len(result.Issues) > 0 {
			debugPrint("First item: #%d - %s", result.Issues[0].GetNumber(), result.Issues[0].GetTitle())
		}

		for _, issue := range result.Issues {
			pr := convertIssueToPR(issue)
			allPRs = append(allPRs, pr)
			if len(allPRs) >= c.pageSize {
				debugPrint("Reached pageSize of %d PRs", c.pageSize)
				break
			}
		}

		if len(allPRs) >= c.pageSize {
			break
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
		page++
	}

	debugPrint("Total PRs fetched: %d of %d", len(allPRs), totalCount)
	return allPRs, totalCount, nil
}

func convertIssueToPR(issue *github.Issue) types.PR {
	var status types.PRStatus
	switch issue.GetState() {
	case "open":
		status = types.StatusOpen
	case "closed":
		if issue.PullRequestLinks != nil && !issue.PullRequestLinks.GetMergedAt().IsZero() {
			status = types.StatusMerged
		} else {
			status = types.StatusClosed
		}
	default:
		status = types.StatusOpen
	}

	var labels []string
	for _, l := range issue.Labels {
		labels = append(labels, l.GetName())
	}

	org, repo := parseRepoName(issue)

	return types.PR{
		Number: issue.GetNumber(),
		Title:  issue.GetTitle(),
		Org:    org,
		Repo:   repo,
		Author: issue.User.GetLogin(),
		Status: status,
		Labels: labels,
		URL:    issue.GetHTMLURL(),
	}
}

func parseRepoName(issue *github.Issue) (org, repo string) {
	repoURL := issue.GetRepositoryURL()
	if repoURL == "" {
		return "", ""
	}

	idx := strings.LastIndex(repoURL, "/repos/")
	if idx == -1 {
		return "", ""
	}
	ownerRepo := repoURL[idx+7:]
	parts := strings.SplitN(ownerRepo, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
