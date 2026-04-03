package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BaconIsAVeg/gh-purview/internal/types"
)

type ReviewDecision string

const (
	ReviewApproved         ReviewDecision = "APPROVED"
	ReviewChangesRequested ReviewDecision = "CHANGES_REQUESTED"
	ReviewRequired         ReviewDecision = "REVIEW_REQUIRED"
)

type graphqlSearchResponse struct {
	Search struct {
		IssueCount int `json:"issueCount"`
		PageInfo   struct {
			HasNextPage bool   `json:"hasNextPage"`
			EndCursor   string `json:"endCursor"`
		} `json:"pageInfo"`
		Nodes []graphqlPRNode `json:"nodes"`
	} `json:"search"`
}

type graphqlPRNode struct {
	Number         int    `json:"number"`
	Title          string `json:"title"`
	State          string `json:"state"`
	ReviewDecision string `json:"reviewDecision"`
	IsDraft        bool   `json:"isDraft"`
	URL            string `json:"url"`
	Author         struct {
		Login string `json:"login"`
	} `json:"author"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Labels struct {
		Nodes []struct {
			Name string `json:"name"`
		} `json:"nodes"`
	} `json:"labels"`
}

func (c *Client) FetchPRsGraphQL(ctx context.Context) ([]types.PR, int, error) {
	query := c.Query()
	debugPrint("GraphQL Query: %s", query)

	var allPRs []types.PR
	var totalCount int
	cursor := ""
	pageSize := c.pageSize
	if pageSize > 100 {
		pageSize = 100
	}

	for {
		prs, total, nextCursor, err := c.fetchPRsPage(ctx, query, pageSize, cursor)
		if err != nil {
			return nil, 0, err
		}

		if totalCount == 0 {
			totalCount = total
		}

		allPRs = append(allPRs, prs...)

		if len(allPRs) >= c.pageSize {
			break
		}
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	debugPrint("Total PRs fetched: %d of %d", len(allPRs), totalCount)
	return allPRs, totalCount, nil
}

func (c *Client) fetchPRsPage(ctx context.Context, query string, first int, after string) ([]types.PR, int, string, error) {
	if c.GraphQL() == nil {
		return nil, 0, "", fmt.Errorf("GraphQL client not initialized")
	}

	graphQLQuery := buildSearchQuery()
	variables := map[string]interface{}{
		"query": query,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	debugPrint("GraphQL query string: %s", graphQLQuery)
	debugPrint("GraphQL variables: %+v", variables)

	var response graphqlSearchResponse
	err := c.GraphQL().Do(graphQLQuery, variables, &response)
	if err != nil {
		debugPrint("GraphQL error: %v", err)
		return nil, 0, "", fmt.Errorf("GraphQL query failed: %w", err)
	}

	respJSON, _ := json.Marshal(response)
	debugPrint("GraphQL response: %s", string(respJSON))

	search := response.Search
	prs := make([]types.PR, 0, len(search.Nodes))

	for _, node := range search.Nodes {
		pr := convertGraphQLNodeToPR(node)
		prs = append(prs, pr)
	}

	nextCursor := ""
	if search.PageInfo.HasNextPage {
		nextCursor = search.PageInfo.EndCursor
	}

	return prs, search.IssueCount, nextCursor, nil
}

func buildSearchQuery() string {
	return `query($query: String!, $first: Int!, $after: String) {
		search(query: $query, type: ISSUE, first: $first, after: $after) {
			issueCount
			pageInfo {
				hasNextPage
				endCursor
			}
			nodes {
				... on PullRequest {
					number
					title
					state
					reviewDecision
					isDraft
					url
					author { login }
					repository {
						name
						owner { login }
					}
					labels(first: 10) {
						nodes { name }
					}
				}
			}
		}
	}`
}

func convertGraphQLNodeToPR(node graphqlPRNode) types.PR {
	var status types.PRStatus
	switch node.State {
	case "OPEN":
		status = types.StatusOpen
	case "CLOSED":
		status = types.StatusClosed
	case "MERGED":
		status = types.StatusMerged
	default:
		status = types.StatusOpen
	}

	var labels []string
	for _, l := range node.Labels.Nodes {
		labels = append(labels, l.Name)
	}

	return types.PR{
		Number:         node.Number,
		Title:          node.Title,
		Org:            node.Repository.Owner.Login,
		Repo:           node.Repository.Name,
		Author:         node.Author.Login,
		Status:         status,
		Labels:         labels,
		URL:            node.URL,
		ReviewDecision: node.ReviewDecision,
	}
}
