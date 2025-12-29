package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListDeploymentPipelines fetches available deployment pipelines for an org
func (c *Client) ListDeploymentPipelines(orgName string) ([]DeploymentPipelineResponse, error) {
	path := "/orgs/" + orgName + "/deployment-pipelines"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp DeploymentPipelineListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.DeploymentPipelines, nil
}
