package github

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/google/go-github/v82/github"
)

type Client struct {
	client        *github.Client
	graphqlClient *api.GraphQLClient
	pageSize      int
	customQuery   string
}

func NewClient(pageSize int) (*Client, error) {
	token, host, err := detectAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to detect auth: %w", err)
	}

	debugPrint("Auth detected: host=%s, token_len=%d", host, len(token))

	httpClient := &http.Client{}
	if token != "" {
		httpClient.Transport = &authTransport{
			token:   token,
			baseURL: host,
		}
	}

	client := github.NewClient(httpClient)
	if host != "" && host != "github.com" {
		client, err = github.NewClient(httpClient).WithEnterpriseURLs(host, host)
		if err != nil {
			return nil, fmt.Errorf("failed to create enterprise client: %w", err)
		}
	}

	debugPrint("Client created, pageSize: %d", pageSize)

	var graphqlClient *api.GraphQLClient
	var gqlErr error
	if token != "" {
		graphqlClient, gqlErr = api.NewGraphQLClient(api.ClientOptions{
			AuthToken:    token,
			LogIgnoreEnv: true,
		})
		debugPrint("GraphQL client created with explicit token")
	} else if host == "github.com" {
		graphqlClient, gqlErr = api.DefaultGraphQLClient()
		debugPrint("GraphQL client using default auth")
	} else {
		graphqlClient, gqlErr = api.NewGraphQLClient(api.ClientOptions{
			Host:         host,
			LogIgnoreEnv: true,
		})
		debugPrint("GraphQL client for enterprise host: %s", host)
	}
	if gqlErr != nil {
		debugPrint("GraphQL client init failed: %v, falling back to REST", gqlErr)
	} else {
		debugPrint("GraphQL client initialized successfully")
	}

	return &Client{
		client:        client,
		graphqlClient: graphqlClient,
		pageSize:      pageSize,
	}, nil
}

func detectAuth() (token, host string, err error) {
	if token = os.Getenv("GH_TOKEN"); token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	debugPrint("Token from env: found=%v", token != "")

	if host = os.Getenv("GH_HOST"); host == "" {
		host, _ = auth.DefaultHost()
	}

	if token == "" {
		token, _ = auth.TokenForHost(host)
		if token == "" {
			return "", "", fmt.Errorf("no authentication token found")
		}
		debugPrint("Using gh auth token for host: %s", host)
	}

	return token, host, nil
}

type authTransport struct {
	token   string
	baseURL string
	base    http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.base == nil {
		t.base = http.DefaultTransport
	}
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "token "+t.token)
	debugPrint("Request: %s %s", req.Method, req.URL.String())
	return t.base.RoundTrip(req)
}

func (c *Client) Close() {}

func (c *Client) Query() string {
	if c.customQuery != "" {
		return c.customQuery
	}
	return prSearchQuery
}

func (c *Client) SetQuery(query string) {
	c.customQuery = query
}
