package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListOrganizations fetches all organizations
func (c *Client) ListOrganizations() ([]OrganizationResponse, error) {
	resp, err := c.doRequest("GET", "/orgs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Always close the body when done!

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Decode JSON response into slice of OrganizationResponse
	var listResp OrganizationListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Organizations, nil
}
