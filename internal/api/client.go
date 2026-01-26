package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// doRequestWithBody performs an HTTP request with a JSON body
func (c *Client) doRequestWithBody(method, path string, body interface{}) (*http.Response, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set(c.APIKeyHeader, c.APIKeyValue)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
}

// doDelete performs a DELETE request and handles common response codes
func (c *Client) doDelete(path string) error {
	resp, err := c.doRequest("DELETE", path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// TestConnection checks if the API server is reachable (without auth)
func TestConnection(baseURL string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	// Try to reach the base URL
	resp, err := client.Get(baseURL)
	if err != nil {
		return fmt.Errorf("cannot connect to server: %w", err)
	}
	defer resp.Body.Close()

	// Any response means server is reachable
	return nil
}

// ValidateAuth tests if credentials are valid and returns organizations if successful
func (c *Client) ValidateAuth() ([]OrganizationResponse, error) {
	orgs, _, err := c.ListOrganizations(DefaultListOptions())
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
