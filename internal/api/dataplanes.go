package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListDataPlanes fetches all data planes for an organization
func (c *Client) ListDataPlanes(orgName string) ([]DataPlane, error) {
	path := "/orgs/" + orgName + "/data-planes"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp DataPlaneListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.DataPlanes, nil
}
