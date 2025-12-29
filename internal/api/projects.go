package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateProject creates a new project in an organization
func (c *Client) CreateProject(orgName string, req CreateProjectRequest) (*ProjectResponse, error) {
	path := "/orgs/" + orgName + "/projects"

	resp, err := c.doRequestWithBody("POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Backend returns 202 Accepted for async operations
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var project ProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &project, nil
}

// ListProjects fetches all projects in an organization
func (c *Client) ListProjects(orgName string) ([]ProjectResponse, error) {
	path := "/orgs/" + orgName + "/projects"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp ProjectListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Projects, nil
}

// GetProject fetches a specific project
func (c *Client) GetProject(orgName, projectName string) (*ProjectResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var project ProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &project, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(orgName, projectName string) error {
	path := "/orgs/" + orgName + "/projects/" + projectName
	return c.doDelete(path)
}
