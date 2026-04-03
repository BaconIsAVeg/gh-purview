package github

import (
	"context"
	"fmt"

	"github.com/BaconIsAVeg/gh-purview/internal/types"
	"github.com/google/go-github/v82/github"
)

func (c *Client) ApprovePR(ctx context.Context, pr *types.PR) error {
	if pr == nil {
		return fmt.Errorf("PR is nil")
	}

	review := &github.PullRequestReviewRequest{
		Event: github.String("APPROVE"),
	}

	_, _, err := c.REST().PullRequests.CreateReview(ctx, pr.Org, pr.Repo, pr.Number, review)
	if err != nil {
		debugPrint("Approve PR error: %v", err)
		return fmt.Errorf("failed to approve PR: %w", err)
	}

	debugPrint("PR #%d approved successfully", pr.Number)
	return nil
}
