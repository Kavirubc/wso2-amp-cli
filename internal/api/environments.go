package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListEnvironments fetches environments for an organization with pagination
func (c *Client) ListEnvironments(orgName string, opts ListOptions) ([]Environment, int, error) {
	params := url.Values{}
	buildPaginationQuery(params, opts)

	path := "/orgs/" + orgName + "/environments"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp EnvironmentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Environments, listResp.Total, nil
}
