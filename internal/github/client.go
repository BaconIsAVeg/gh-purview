package github

import (
	"github.com/BaconIsAVeg/github-tuis/debug"
	tuisclient "github.com/BaconIsAVeg/github-tuis/github/client"
)

type Client struct {
	*tuisclient.BaseClient
	pageSize    int
	customQuery string
}

func NewClient(pageSize int) (*Client, error) {
	bc, err := tuisclient.NewClient(tuisclient.ClientOptions{})
	if err != nil {
		return nil, err
	}

	debug.Print("Client created, pageSize: %d", pageSize)

	return &Client{
		BaseClient: bc,
		pageSize:   pageSize,
	}, nil
}

func (c *Client) Query() string {
	if c.customQuery != "" {
		return c.customQuery
	}
	return prSearchQuery
}

func (c *Client) SetQuery(query string) {
	c.customQuery = query
}
