package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListEnvironments fetches all environments for an organization
func (c *Client) ListEnvironments(orgName string) ([]Environment, error) {
	path := "/orgs/" + orgName + "/environments"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp EnvironmentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Environments, nil
}
