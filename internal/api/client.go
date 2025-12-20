package api

import (
	"net/http"
	"time"
)

// Client is the HTTP client for the Agent Management Platform API
type Client struct {
	BaseURL      string
	APIKeyHeader string
	APIKeyValue  string
	HTTPClient   *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, apiKeyHeader, apiKeyValue string) *Client {
	return &Client{
		BaseURL:      baseURL,
		APIKeyHeader: apiKeyHeader,
		APIKeyValue:  apiKeyValue,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string) (*http.Response, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set authentication and content type headers
	req.Header.Set(c.APIKeyHeader, c.APIKeyValue)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
}
