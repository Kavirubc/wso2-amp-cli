package api

import (
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

// ==================== Organizations ====================

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

// ==================== Projects ====================

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

// ==================== Agents ====================

// ListAgents fetches all agents in a project
func (c *Client) ListAgents(orgName, projectName string) ([]AgentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp AgentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Agents, nil
}

// GetAgent fetches a specific agent
func (c *Client) GetAgent(orgName, projectName, agentName string) (*AgentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var agent AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &agent, nil // Return pointer to the agent
}

// ==================== Builds ====================

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

// ==================== Deployments ====================

// ListDeployments fetches all deployments for an agent
func (c *Client) ListDeployments(orgName, projectName, agentName string) ([]DeploymentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/deployments"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var deployments []DeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return deployments, nil
}
