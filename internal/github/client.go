package github

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/google/go-github/v82/github"
)

type Client struct {
	client   *github.Client
	org      string
	pageSize int
}

func NewClient(org string, pageSize int) (*Client, error) {
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

	debugPrint("Client created for org: %s, pageSize: %d", org, pageSize)

	return &Client{
		client:   client,
		org:      org,
		pageSize: pageSize,
	}, nil
}

func detectAuth() (token, host string, err error) {
	if token = os.Getenv("GH_TOKEN"); token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	debugPrint("Token from env: found=%v", token != "")

	if host = os.Getenv("GH_HOST"); host == "" {
		host = "github.com"
	}

	if token == "" {
		restClient, err := api.DefaultRESTClient()
		if err != nil {
			return "", "", fmt.Errorf("no token found and gh auth unavailable: %w", err)
		}
		_ = restClient
		debugPrint("Using gh auth fallback")
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
	return fmt.Sprintf(prSearchQuery, c.org)
}
