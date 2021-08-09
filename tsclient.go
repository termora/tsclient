// Package tsclient is a client library for the Typesense search engine.
package tsclient

import (
	"encoding/json"
	"net/http"
	"strings"
)

// VERSION is unlikely to ever be updated even as the library gets new releases
const VERSION = "0.1.0"

// Client is a Typesense API client.
type Client struct {
	Client *http.Client

	baseURL string
	apiKey  string

	// Debug is a debug logging function. No-op by default.
	Debug func(tmpl string, args ...interface{})

	UserAgent string
}

// New creates a new Client and pings the server.
func New(url, apiKey string) (*Client, error) {
	c := &Client{
		Client:    &http.Client{},
		baseURL:   strings.TrimSuffix(url, "/"),
		apiKey:    apiKey,
		Debug:     func(string, ...interface{}) {},
		UserAgent: "go/tsclient " + VERSION,
	}

	// we only care about if the request goes through at all
	_, err := c.Health()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Health performs a health check on the server.
func (c *Client) Health() (ok bool, err error) {
	resp, err := c.Request("GET", "/health")
	if err != nil {
		return
	}

	s := struct {
		OK bool `json:"ok"`
	}{}

	err = json.Unmarshal(resp, &s)
	if err != nil {
		return
	}

	return s.OK, nil
}
