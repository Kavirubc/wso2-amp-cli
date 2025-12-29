package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TriggerBuild triggers a new build for an agent
func (c *Client) TriggerBuild(orgName, projectName, agentName, commitID string) (*BuildResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/builds"
	if commitID != "" {
		path += "?commitId=" + url.QueryEscape(commitID)
	}

	resp, err := c.doRequestWithBody("POST", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var build BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &build, nil
}

// ListBuilds fetches all builds for an agent
func (c *Client) ListBuilds(orgName, projectName, agentName string) ([]BuildResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/builds"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var builds []BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return builds, nil
}

// GetBuild fetches a specific build
func (c *Client) GetBuild(orgName, projectName, agentName, buildName string) (*BuildResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/builds/" + buildName

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var build BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &build, nil
}
